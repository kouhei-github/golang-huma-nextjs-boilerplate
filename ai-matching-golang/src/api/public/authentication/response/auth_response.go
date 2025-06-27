package response

import (
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	AccessToken          string    `json:"accessToken,omitempty" doc:"JWT access token"`
	IdToken              string    `json:"idToken,omitempty" doc:"JWT ID token"`
	RefreshToken         string    `json:"refreshToken,omitempty" doc:"JWT refresh token"`
	TokenType            string    `json:"tokenType,omitempty" doc:"Token type (Bearer)"`
	ExpiresAt            time.Time `json:"expiresAt,omitempty" doc:"Token expiration time"`
	User                 UserInfo  `json:"user" doc:"User information"`
	Message              string    `json:"message,omitempty" doc:"Response message"`
	RequiresConfirmation bool      `json:"requiresConfirmation,omitempty" doc:"Whether email confirmation is required"`
}

type UserInfo struct {
	ID        uuid.UUID `json:"id" doc:"User ID"`
	Email     string    `json:"email" doc:"User email"`
	FirstName string    `json:"firstName" doc:"User first name"`
	LastName  string    `json:"lastName" doc:"User last name"`
	Tenants   []string  `json:"tenants,omitempty" doc:"List of tenant names user belongs to"`
}
