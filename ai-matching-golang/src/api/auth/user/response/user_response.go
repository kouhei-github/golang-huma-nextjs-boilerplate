package response

import "time"

type UserResponse struct {
	ID             int64     `json:"id" doc:"User ID"`
	Email          string    `json:"email" doc:"User email"`
	FirstName      string    `json:"firstName" doc:"User first name"`
	LastName       string    `json:"lastName" doc:"User last name"`
	OrganizationID *int64    `json:"organizationId,omitempty" doc:"Organization ID"`
	TenantID       *int64    `json:"tenantId,omitempty" doc:"Tenant ID"`
	CreatedAt      time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt      time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

type UserListResponse struct {
	Users    []UserResponse `json:"users" doc:"List of users"`
	Total    int            `json:"total" doc:"Total count"`
	Page     int            `json:"page" doc:"Current page"`
	PageSize int            `json:"pageSize" doc:"Page size"`
}