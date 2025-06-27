package requests

import "github.com/google/uuid"

type CreateTenantRequest struct {
	OrganizationID uuid.UUID `json:"organizationId" validate:"required" doc:"Organization ID"`
	Name           string `json:"name" validate:"required" doc:"Tenant name"`
	Subdomain      string `json:"subdomain" validate:"required" doc:"Tenant subdomain"`
	IsActive       bool   `json:"isActive" doc:"Is tenant active"`
}

type UpdateTenantRequest struct {
	Name      string `json:"name" validate:"required" doc:"Tenant name"`
	Subdomain string `json:"subdomain" validate:"required" doc:"Tenant subdomain"`
	IsActive  bool   `json:"isActive" doc:"Is tenant active"`
}