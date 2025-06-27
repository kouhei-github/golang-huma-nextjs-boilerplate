package requests

type CreateUserRequest struct {
	Email      string  `json:"email" validate:"required,email" doc:"User email"`
	CognitoID  string  `json:"cognitoId" validate:"required" doc:"Cognito user ID"`
	FirstName  string  `json:"firstName" validate:"required" doc:"User first name"`
	LastName   string  `json:"lastName" validate:"required" doc:"User last name"`
	TenantID   *int64  `json:"tenantId,omitempty" doc:"Initial tenant ID to associate user with"`
	TenantRole *string `json:"tenantRole,omitempty" doc:"Role in the initial tenant"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"required,email" doc:"User email"`
	FirstName string `json:"firstName" validate:"required" doc:"User first name"`
	LastName  string `json:"lastName" validate:"required" doc:"User last name"`
}