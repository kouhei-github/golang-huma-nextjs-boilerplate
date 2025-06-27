package middleware

import (
	"context"
	"fmt"
	"strings"

	"ai-matching/src/infrastructure/external/cognito"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtValidator *cognito.CognitoJWTValidator
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtValidator: cognito.NewCognitoJWTValidator(),
	}
}

func (m *AuthMiddleware) FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token, err := m.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: " + err.Error(),
			})
		}

		userInfo, err := m.jwtValidator.GetUserInfoFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to extract user info: " + err.Error(),
			})
		}

		c.Locals("user_id", userInfo["user_id"])
		c.Locals("email", userInfo["email"])
		c.Locals("token", tokenString)

		return c.Next()
	}
}

func (m *AuthMiddleware) HumaMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			ctx.SetStatus(401)
			ctx.SetHeader("Content-Type", "application/json")
			ctx.BodyWriter().Write([]byte(`{"title":"Unauthorized","status":401,"detail":"Missing authorization header"}`))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.SetStatus(401)
			ctx.SetHeader("Content-Type", "application/json")
			ctx.BodyWriter().Write([]byte(`{"title":"Unauthorized","status":401,"detail":"Invalid authorization header format"}`))
			return
		}

		token, err := m.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			ctx.SetStatus(401)
			ctx.SetHeader("Content-Type", "application/json")
			detail := fmt.Sprintf(`{"title":"Unauthorized","status":401,"detail":"Invalid token: %s"}`, err.Error())
			ctx.BodyWriter().Write([]byte(detail))
			return
		}

		userInfo, err := m.jwtValidator.GetUserInfoFromToken(token)
		if err != nil {
			ctx.SetStatus(401)
			ctx.SetHeader("Content-Type", "application/json")
			detail := fmt.Sprintf(`{"title":"Unauthorized","status":401,"detail":"Failed to extract user info: %s"}`, err.Error())
			ctx.BodyWriter().Write([]byte(detail))
			return
		}

		newCtx := context.WithValue(ctx.Context(), "user_id", userInfo["user_id"])
		newCtx = context.WithValue(newCtx, "email", userInfo["email"])
		newCtx = context.WithValue(newCtx, "token", tokenString)

		ctx = huma.WithContext(ctx, newCtx)

		next(ctx)
	}
}

type UserContext struct {
	UserID string
	Email  string
	Token  string
}

func GetUserFromContext(ctx context.Context) (*UserContext, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("user_id not found in context")
	}

	email, ok := ctx.Value("email").(string)
	if !ok {
		return nil, fmt.Errorf("email not found in context")
	}

	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	return &UserContext{
		UserID: userID,
		Email:  email,
		Token:  token,
	}, nil
}
