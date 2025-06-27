package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/public/authentication/requests"
	"ai-matching/src/api/public/authentication/response"
	"ai-matching/src/domain/interface/external"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	userRepo         repository.UserRepository
	tenantRepo       repository.TenantRepository
	tenantUserRepo   repository.TenantUserRepository
	organizationRepo repository.OrganizationRepository
	cognitoClient    external.CognitoClient
}

func NewAuthUsecase(userRepo repository.UserRepository, tenantUserRepo repository.TenantUserRepository, tenantRepo repository.TenantRepository, organizationRepo repository.OrganizationRepository, cognitoClient external.CognitoClient) *AuthUsecase {
	return &AuthUsecase{
		userRepo:         userRepo,
		tenantUserRepo:   tenantUserRepo,
		tenantRepo:       tenantRepo,
		organizationRepo: organizationRepo,
		cognitoClient:    cognitoClient,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, req requests.LoginRequest) (*response.AuthResponse, error) {
	authResult, err := u.cognitoClient.InitiateAuth(ctx, req.Email, req.Password)
	if err != nil {
		if awsErr, ok := err.(*types.NotAuthorizedException); ok {
			_ = awsErr
			return nil, errors.New("invalid credentials")
		}
		if awsErr, ok := err.(*types.UserNotFoundException); ok {
			_ = awsErr
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if authResult.AuthenticationResult == nil {
		return nil, errors.New("authentication failed: no result")
	}

	// Extract tokens
	idToken := aws.ToString(authResult.AuthenticationResult.IdToken)

	// Extract sub claim from ID token
	cognitoUserID, err := extractSubFromIDToken(idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract user ID from token: %w", err)
	}

	// Try to get user by cognito ID first, then by email
	user, err := u.userRepo.GetUserByCognitoID(ctx, cognitoUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User exists in Cognito but not in local DB, try by email
			user, err = u.userRepo.GetUserByEmail(ctx, req.Email)
			if err != nil {
				if err == sql.ErrNoRows {
					// Create user in local DB since they exist in Cognito
					user, err = u.userRepo.CreateUser(ctx, db.CreateUserParams{
						CognitoID: cognitoUserID,
						Email:     req.Email,
						FirstName: sql.NullString{}, // Will be populated from Cognito claims if needed
						LastName:  sql.NullString{},
					})
					if err != nil {
						return nil, fmt.Errorf("failed to create user in local database: %w", err)
					}
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	accessToken := aws.ToString(authResult.AuthenticationResult.AccessToken)
	refreshToken := aws.ToString(authResult.AuthenticationResult.RefreshToken)
	expiresIn := authResult.AuthenticationResult.ExpiresIn
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Refresh token is managed by Cognito, no need to store locally

	return &response.AuthResponse{
		AccessToken:  accessToken,
		IdToken:      idToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User: response.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
			Tenants:   []string{}, // TODO: Fetch user's tenants
		},
	}, nil
}

func (u *AuthUsecase) Register(ctx context.Context, req requests.RegisterRequest) (*response.AuthResponse, error) {
	if req.OrganizationName == nil || req.TenantName == nil || req.TenantSubdomain == nil {
		return nil, errors.New("either complete organization/tenant information must be provided")
	}

	attributes := map[string]string{
		"email":       req.Email,
		"given_name":  req.FirstName,
		"family_name": req.LastName,
	}

	signUpResult, err := u.cognitoClient.SignUp(ctx, req.Email, req.Password, attributes)
	if err != nil {
		if awsErr, ok := err.(*types.UsernameExistsException); ok {
			_ = awsErr
			return nil, errors.New("user already exists")
		}
		if awsErr, ok := err.(*types.InvalidPasswordException); ok {
			_ = awsErr
			return nil, errors.New("password does not meet requirements")
		}
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	cognitoUserID := aws.ToString(signUpResult.UserSub)

	user, err := u.userRepo.CreateUser(ctx, db.CreateUserParams{
		CognitoID: cognitoUserID,
		Email:     req.Email,
		FirstName: sql.NullString{String: req.FirstName, Valid: true},
		LastName:  sql.NullString{String: req.LastName, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	org, err := u.organizationRepo.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:        *req.OrganizationName,
		Description: sql.NullString{String: aws.ToString(req.OrganizationDescription), Valid: req.OrganizationDescription != nil},
		IsActive:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	tenant, err := u.tenantRepo.CreateTenant(ctx, db.CreateTenantParams{
		OrganizationID: org.ID,
		Name:           *req.TenantName,
		Subdomain:      *req.TenantSubdomain,
		IsActive:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	if _, err := u.tenantUserRepo.AddUserToTenant(ctx, db.AddUserToTenantParams{
		TenantID: tenant.ID,
		UserID:   user.ID,
	}); err != nil {
		return nil, fmt.Errorf("failed to associate user with tenant: %w", err)
	}

	if signUpResult.UserConfirmed {
		return u.Login(ctx, requests.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
	}

	return &response.AuthResponse{
		Message:              "User registered successfully. Please check your email for confirmation code.",
		RequiresConfirmation: true,
		User: response.UserInfo{
			ID:        user.ID,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Tenants:   []string{}, // User not yet assigned to any tenant
		},
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, req requests.RefreshTokenRequest) (*response.AuthResponse, error) {
	// Refresh token validation is handled by Cognito

	authResult, err := u.cognitoClient.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if awsErr, ok := err.(*types.NotAuthorizedException); ok {
			_ = awsErr
			return nil, errors.New("refresh token expired or invalid")
		}
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	if authResult.AuthenticationResult == nil {
		return nil, errors.New("token refresh failed: no result")
	}

	// Extract user ID from refresh token to get user details
	// This would typically be done by parsing the token or getting it from Cognito
	// For now, we'll skip user info in refresh response

	accessToken := aws.ToString(authResult.AuthenticationResult.AccessToken)
	idToken := aws.ToString(authResult.AuthenticationResult.IdToken)
	newRefreshToken := req.RefreshToken
	if authResult.AuthenticationResult.RefreshToken != nil {
		newRefreshToken = aws.ToString(authResult.AuthenticationResult.RefreshToken)
	}
	expiresIn := authResult.AuthenticationResult.ExpiresIn
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Refresh token is managed by Cognito, no need to store locally

	return &response.AuthResponse{
		AccessToken:  accessToken,
		IdToken:      idToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		// User info is not included in refresh token response
	}, nil
}

func (u *AuthUsecase) ConfirmSignUp(ctx context.Context, email, confirmationCode string) error {
	err := u.cognitoClient.ConfirmSignUp(ctx, email, confirmationCode)
	if err != nil {
		if awsErr, ok := err.(*types.CodeMismatchException); ok {
			_ = awsErr
			return errors.New("invalid confirmation code")
		}
		if awsErr, ok := err.(*types.ExpiredCodeException); ok {
			_ = awsErr
			return errors.New("confirmation code has expired")
		}
		return fmt.Errorf("confirmation failed: %w", err)
	}
	return nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	_, err := u.cognitoClient.ForgotPassword(ctx, email)
	if err != nil {
		if awsErr, ok := err.(*types.UserNotFoundException); ok {
			_ = awsErr
			return nil
		}
		return fmt.Errorf("forgot password failed: %w", err)
	}
	return nil
}

func (u *AuthUsecase) ConfirmForgotPassword(ctx context.Context, email, password, confirmationCode string) error {
	err := u.cognitoClient.ConfirmForgotPassword(ctx, email, password, confirmationCode)
	if err != nil {
		if awsErr, ok := err.(*types.CodeMismatchException); ok {
			_ = awsErr
			return errors.New("invalid confirmation code")
		}
		if awsErr, ok := err.(*types.ExpiredCodeException); ok {
			_ = awsErr
			return errors.New("confirmation code has expired")
		}
		if awsErr, ok := err.(*types.InvalidPasswordException); ok {
			_ = awsErr
			return errors.New("password does not meet requirements")
		}
		return fmt.Errorf("password reset failed: %w", err)
	}

	// Password is managed by Cognito, no need to update locally

	return nil
}

func nullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func ptrToNullInt64(p *int64) sql.NullInt64 {
	if p != nil {
		return sql.NullInt64{Int64: *p, Valid: true}
	}
	return sql.NullInt64{}
}

// extractSubFromIDToken extracts the sub claim (Cognito User ID) from the ID token
func extractSubFromIDToken(idToken string) (string, error) {
	// Parse the token without validation (validation is done by Cognito)
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse ID token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errors.New("sub claim not found in ID token")
	}

	return sub, nil
}
