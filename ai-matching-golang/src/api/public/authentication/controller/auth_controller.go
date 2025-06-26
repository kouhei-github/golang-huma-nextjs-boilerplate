package controller

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/public/authentication/requests"
	"ai-matching/src/api/public/authentication/response"
	"ai-matching/src/api/public/authentication/usecase"
	"ai-matching/src/infrastructure/repository"
	"context"
)

type AuthController struct {
	usecase *usecase.AuthUsecase
}

func NewAuthController(queries db.Querier) *AuthController {
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)
	authUsecase := usecase.NewAuthUsecase(userRepo, authRepo)

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
