package external

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoClient interface {
	SignUp(ctx context.Context, email, password string, attributes map[string]string) (*cognitoidentityprovider.SignUpOutput, error)
	ConfirmSignUp(ctx context.Context, email, confirmationCode string) error
	InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error)
	ForgotPassword(ctx context.Context, email string) (*cognitoidentityprovider.ForgotPasswordOutput, error)
	ConfirmForgotPassword(ctx context.Context, email, password, confirmationCode string) error
	GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error)
}
