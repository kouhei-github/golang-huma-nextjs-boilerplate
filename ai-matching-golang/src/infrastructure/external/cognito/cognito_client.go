package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoClient struct {
	client       *cognitoidentityprovider.Client
	userPoolID   string
	clientID     string
	clientSecret string
}

func NewCognitoClient() (*CognitoClient, error) {
	cognitoEndpoint := os.Getenv("COGNITO_ENDPOINT")
	if cognitoEndpoint == "" {
		cognitoEndpoint = "https://cognito-idp.us-east-1.amazonaws.com"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == cognitoidentityprovider.ServiceID && cognitoEndpoint != "" {
					return aws.Endpoint{
						URL:               cognitoEndpoint,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &CognitoClient{
		client:       cognitoidentityprovider.NewFromConfig(cfg),
		userPoolID:   os.Getenv("COGNITO_USER_POOL_ID"),
		clientID:     os.Getenv("COGNITO_CLIENT_ID"),
		clientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
	}, nil
}

func (c *CognitoClient) SignUp(ctx context.Context, email, password string, attributes map[string]string) (*cognitoidentityprovider.SignUpOutput, error) {
	secretHash := c.calculateSecretHash(email)

	userAttributes := []types.AttributeType{
		{
			Name:  aws.String("email"),
			Value: aws.String(email),
		},
	}

	for key, value := range attributes {
		userAttributes = append(userAttributes, types.AttributeType{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(c.clientID),
		Username:       aws.String(email),
		Password:       aws.String(password),
		SecretHash:     aws.String(secretHash),
		UserAttributes: userAttributes,
	}

	result, err := c.client.SignUp(ctx, input)
	if err != nil {
		return nil, err
	}

	// Auto-confirm the user for development environment
	if os.Getenv("COGNITO_AUTO_CONFIRM") == "true" {
		confirmInput := &cognitoidentityprovider.AdminConfirmSignUpInput{
			UserPoolId: aws.String(c.userPoolID),
			Username:   aws.String(email),
		}
		_, confirmErr := c.client.AdminConfirmSignUp(ctx, confirmInput)
		if confirmErr != nil {
			// Log the error but don't fail the signup
			fmt.Printf("Warning: Failed to auto-confirm user %s: %v\n", email, confirmErr)
		} else {
			// Mark the user as confirmed in the result
			result.UserConfirmed = true
		}
	}

	return result, nil
}

func (c *CognitoClient) ConfirmSignUp(ctx context.Context, email, confirmationCode string) error {
	secretHash := c.calculateSecretHash(email)

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(secretHash),
	}

	_, err := c.client.ConfirmSignUp(ctx, input)
	return err
}

func (c *CognitoClient) InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	secretHash := c.calculateSecretHash(email)
	authParams := map[string]string{
		"USERNAME":    email,
		"PASSWORD":    password,
		"SECRET_HASH": secretHash,
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       aws.String(c.clientID),
		AuthParameters: authParams,
	}

	return c.client.InitiateAuth(ctx, input)
}

func (c *CognitoClient) RefreshToken(ctx context.Context, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	authParams := map[string]string{
		"REFRESH_TOKEN": refreshToken,
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeRefreshTokenAuth,
		ClientId:       aws.String(c.clientID),
		AuthParameters: authParams,
	}

	return c.client.InitiateAuth(ctx, input)
}

func (c *CognitoClient) ForgotPassword(ctx context.Context, email string) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	secretHash := c.calculateSecretHash(email)

	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   aws.String(c.clientID),
		Username:   aws.String(email),
		SecretHash: aws.String(secretHash),
	}

	return c.client.ForgotPassword(ctx, input)
}

func (c *CognitoClient) ConfirmForgotPassword(ctx context.Context, email, password, confirmationCode string) error {
	secretHash := c.calculateSecretHash(email)

	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(email),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(secretHash),
	}

	_, err := c.client.ConfirmForgotPassword(ctx, input)
	return err
}

func (c *CognitoClient) GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	return c.client.GetUser(ctx, input)
}

func (c *CognitoClient) calculateSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(c.clientSecret))
	mac.Write([]byte(username + c.clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
