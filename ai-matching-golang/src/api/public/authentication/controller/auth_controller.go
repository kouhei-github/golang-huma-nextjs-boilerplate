package controller

import (
	"ai-matching/src/api/public/authentication/requests"
	"ai-matching/src/api/public/authentication/response"
	"ai-matching/src/api/public/authentication/usecase"
	"context"
)

type AuthController struct {
	usecase *usecase.AuthUsecase
}

func NewAuthController(authUsecase *usecase.AuthUsecase) *AuthController {
	return &AuthController{
		usecase: authUsecase,
	}
}

type LoginInput struct {
	Body requests.LoginRequest
}

type LoginOutput struct {
	Body response.AuthResponse
}

func (c *AuthController) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	resp, err := c.usecase.Login(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Body: *resp}, nil
}

type RegisterInput struct {
	Body requests.RegisterRequest
}

type RegisterOutput struct {
	Body response.AuthResponse
}

func (c *AuthController) Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	resp, err := c.usecase.Register(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{Body: *resp}, nil
}

type RefreshTokenInput struct {
	Body requests.RefreshTokenRequest
}

type RefreshTokenOutput struct {
	Body response.AuthResponse
}

func (c *AuthController) RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
	resp, err := c.usecase.RefreshToken(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenOutput{Body: *resp}, nil
}

type ConfirmSignUpInput struct {
	Body requests.ConfirmSignUpRequest
}

type ConfirmSignUpOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (c *AuthController) ConfirmSignUp(ctx context.Context, input *ConfirmSignUpInput) (*ConfirmSignUpOutput, error) {
	err := c.usecase.ConfirmSignUp(ctx, input.Body.Email, input.Body.ConfirmationCode)
	if err != nil {
		return nil, err
	}

	return &ConfirmSignUpOutput{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "Email confirmed successfully",
		},
	}, nil
}

type ForgotPasswordInput struct {
	Body requests.ForgotPasswordRequest
}

type ForgotPasswordOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (c *AuthController) ForgotPassword(ctx context.Context, input *ForgotPasswordInput) (*ForgotPasswordOutput, error) {
	err := c.usecase.ForgotPassword(ctx, input.Body.Email)
	if err != nil {
		return nil, err
	}

	return &ForgotPasswordOutput{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "If the email exists, a password reset code has been sent",
		},
	}, nil
}

type ConfirmForgotPasswordInput struct {
	Body requests.ConfirmForgotPasswordRequest
}

type ConfirmForgotPasswordOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (c *AuthController) ConfirmForgotPassword(ctx context.Context, input *ConfirmForgotPasswordInput) (*ConfirmForgotPasswordOutput, error) {
	err := c.usecase.ConfirmForgotPassword(ctx, input.Body.Email, input.Body.Password, input.Body.ConfirmationCode)
	if err != nil {
		return nil, err
	}

	return &ConfirmForgotPasswordOutput{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "Password reset successfully",
		},
	}, nil
}
