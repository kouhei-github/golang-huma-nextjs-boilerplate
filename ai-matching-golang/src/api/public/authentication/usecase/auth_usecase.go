package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/public/authentication/requests"
	"ai-matching/src/api/public/authentication/response"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewAuthUsecase(userRepo repository.UserRepository, authRepo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, req requests.LoginRequest) (*response.AuthResponse, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken := "mock_access_token_" + user.Email
	refreshToken := "mock_refresh_token_" + user.Email
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = u.authRepo.CreateUserAuth(ctx, db.CreateUserAuthParams{
		UserID:       user.ID,
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		ExpiresAt:    sql.NullTime{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: response.UserInfo{
			ID:             user.ID,
			Email:          user.Email,
			FirstName:      user.FirstName.String,
			LastName:       user.LastName.String,
			OrganizationID: nullInt64ToPtr(user.OrganizationID),
			TenantID:       nullInt64ToPtr(user.TenantID),
		},
	}, nil
}

func (u *AuthUsecase) Register(ctx context.Context, req requests.RegisterRequest) (*response.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = u.userRepo.CreateUser(ctx, db.CreateUserParams{
		Email:          req.Email,
		PasswordHash:   string(hashedPassword),
		FirstName:      sql.NullString{String: req.FirstName, Valid: true},
		LastName:       sql.NullString{String: req.LastName, Valid: true},
		OrganizationID: ptrToNullInt64(req.OrgID),
		TenantID:       ptrToNullInt64(req.TenantID),
	})
	if err != nil {
		return nil, err
	}

	return u.Login(ctx, requests.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, req requests.RefreshTokenRequest) (*response.AuthResponse, error) {
	userAuth, err := u.authRepo.GetUserAuthByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := u.userRepo.GetUser(ctx, userAuth.UserID)
	if err != nil {
		return nil, err
	}

	accessToken := "mock_access_token_" + user.Email
	refreshToken := "mock_refresh_token_" + user.Email
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = u.authRepo.UpdateUserAuth(ctx, db.UpdateUserAuthParams{
		ID:           userAuth.ID,
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		ExpiresAt:    sql.NullTime{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: response.UserInfo{
			ID:             user.ID,
			Email:          user.Email,
			FirstName:      user.FirstName.String,
			LastName:       user.LastName.String,
			OrganizationID: nullInt64ToPtr(user.OrganizationID),
			TenantID:       nullInt64ToPtr(user.TenantID),
		},
	}, nil
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
