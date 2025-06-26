package response

import "time"

type TenantResponse struct {
	ID             int64     `json:"id" doc:"Tenant ID"`
	OrganizationID int64     `json:"organizationId" doc:"Organization ID"`
	Name           string    `json:"name" doc:"Tenant name"`
	Subdomain      string    `json:"subdomain" doc:"Tenant subdomain"`
	IsActive       bool      `json:"isActive" doc:"Is tenant active"`
	CreatedAt      time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt      time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

type TenantListResponse struct {
	Tenants  []TenantResponse `json:"tenants" doc:"List of tenants"`
	Total    int              `json:"total" doc:"Total count"`
	Page     int              `json:"page" doc:"Current page"`
	PageSize int              `json:"pageSize" doc:"Page size"`
}