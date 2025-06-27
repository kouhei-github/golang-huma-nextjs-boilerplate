package requests

type AddUserToTenantRequest struct {
	UserID int64  `json:"userId" validate:"required" doc:"User ID to add"`
	Role   string `json:"role" default:"member" doc:"Role of the user in the tenant"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required" doc:"New role for the user"`
}