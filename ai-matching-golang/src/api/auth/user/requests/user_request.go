package requests

type CreateUserRequest struct {
	Email          string  `json:"email" validate:"required,email" doc:"User email"`
	Password       string  `json:"password" validate:"required,min=6" doc:"User password"`
	FirstName      string  `json:"firstName" validate:"required" doc:"User first name"`
	LastName       string  `json:"lastName" validate:"required" doc:"User last name"`
	OrganizationID *int64  `json:"organizationId,omitempty" doc:"Organization ID"`
	TenantID       *int64  `json:"tenantId,omitempty" doc:"Tenant ID"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"required,email" doc:"User email"`
	FirstName string `json:"firstName" validate:"required" doc:"User first name"`
	LastName  string `json:"lastName" validate:"required" doc:"User last name"`
}