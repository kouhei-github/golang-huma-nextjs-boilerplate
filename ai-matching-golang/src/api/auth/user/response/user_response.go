package response

import "time"

type UserResponse struct {
	ID        int64        `json:"id" doc:"User ID"`
	Email     string       `json:"email" doc:"User email"`
	FirstName string       `json:"firstName" doc:"User first name"`
	LastName  string       `json:"lastName" doc:"User last name"`
	Tenants   []TenantInfo `json:"tenants,omitempty" doc:"List of tenants user belongs to"`
	CreatedAt time.Time    `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time    `json:"updatedAt" doc:"Last update timestamp"`
}

type TenantInfo struct {
	ID        int64  `json:"id" doc:"Tenant ID"`
	Name      string `json:"name" doc:"Tenant name"`
	Subdomain string `json:"subdomain" doc:"Tenant subdomain"`
	Role      string `json:"role" doc:"User role in tenant"`
}

type UserListResponse struct {
	Users    []UserResponse `json:"users" doc:"List of users"`
	Total    int            `json:"total" doc:"Total count"`
	Page     int            `json:"page" doc:"Current page"`
	PageSize int            `json:"pageSize" doc:"Page size"`
}