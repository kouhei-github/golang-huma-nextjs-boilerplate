package response

import "time"

type AuthResponse struct {
	AccessToken  string    `json:"accessToken" doc:"JWT access token"`
	RefreshToken string    `json:"refreshToken" doc:"JWT refresh token"`
	ExpiresAt    time.Time `json:"expiresAt" doc:"Token expiration time"`
	User         UserInfo  `json:"user" doc:"User information"`
}

type UserInfo struct {
	ID             int64   `json:"id" doc:"User ID"`
	Email          string  `json:"email" doc:"User email"`
	FirstName      string  `json:"firstName" doc:"User first name"`
	LastName       string  `json:"lastName" doc:"User last name"`
	OrganizationID *int64  `json:"organizationId,omitempty" doc:"Organization ID"`
	TenantID       *int64  `json:"tenantId,omitempty" doc:"Tenant ID"`
}