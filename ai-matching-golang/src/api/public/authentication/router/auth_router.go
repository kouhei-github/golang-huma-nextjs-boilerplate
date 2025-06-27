package router

import (
	"ai-matching/src/api/public/authentication/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(api huma.API, router fiber.Router, authController *controller.AuthController) {

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      "POST",
		Path:        "/api/v1/public/auth/login",
		Summary:     "User login",
		Description: "Authenticate user and get access token",
		Tags:        []string{"Authentication"},
	}, authController.Login)

	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      "POST",
		Path:        "/api/v1/public/auth/register",
		Summary:     "User registration",
		Description: "Register a new user",
		Tags:        []string{"Authentication"},
	}, authController.Register)

	huma.Register(api, huma.Operation{
		OperationID: "refresh-token",
		Method:      "POST",
		Path:        "/api/v1/public/auth/refresh",
		Summary:     "Refresh token",
		Description: "Get new access token using refresh token",
		Tags:        []string{"Authentication"},
	}, authController.RefreshToken)

	huma.Register(api, huma.Operation{
		OperationID: "confirm-signup",
		Method:      "POST",
		Path:        "/api/v1/public/auth/confirm",
		Summary:     "Confirm user registration",
		Description: "Confirm user email with confirmation code",
		Tags:        []string{"Authentication"},
	}, authController.ConfirmSignUp)

	huma.Register(api, huma.Operation{
		OperationID: "forgot-password",
		Method:      "POST",
		Path:        "/api/v1/public/auth/forgot-password",
		Summary:     "Forgot password",
		Description: "Request password reset code",
		Tags:        []string{"Authentication"},
	}, authController.ForgotPassword)

	huma.Register(api, huma.Operation{
		OperationID: "confirm-forgot-password",
		Method:      "POST",
		Path:        "/api/v1/public/auth/reset-password",
		Summary:     "Reset password",
		Description: "Reset password with confirmation code",
		Tags:        []string{"Authentication"},
	}, authController.ConfirmForgotPassword)
}
