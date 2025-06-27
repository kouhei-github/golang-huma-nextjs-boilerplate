package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" doc:"User email"`
	Password string `json:"password" validate:"required,min=6" doc:"User password"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email" doc:"User email"`
	Password  string `json:"password" validate:"required,min=6" doc:"User password"`
	FirstName string `json:"firstName" validate:"required" doc:"User first name"`
	LastName  string `json:"lastName" validate:"required" doc:"User last name"`

	// Organization creation fields (optional - used when creating new organization)
	OrganizationName        *string `json:"organizationName,omitempty" validate:"omitempty,min=1,max=255" doc:"Organization name for new organization"`
	OrganizationDescription *string `json:"organizationDescription,omitempty" doc:"Organization description"`

	// Tenant creation fields (optional - used when creating new tenant)
	TenantName      *string `json:"tenantName,omitempty" validate:"omitempty,min=1,max=255" doc:"Tenant name for new tenant"`
	TenantSubdomain *string `json:"tenantSubdomain,omitempty" validate:"omitempty,min=1,max=100,alphanum" doc:"Unique subdomain for new tenant"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" doc:"Refresh token"`
}
