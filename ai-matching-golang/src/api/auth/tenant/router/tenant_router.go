package router

import (
	"ai-matching/src/api/auth/tenant/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterTenantRoutes(api huma.API, router fiber.Router, tenantController *controller.TenantController) {
	// Organization-scoped tenant endpoints
	huma.Register(api, huma.Operation{
		OperationID: "list-tenants-by-organization",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/tenants",
		Summary:     "List tenants by organization",
		Description: "List all tenants for an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.ListTenantsByOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "get-tenant-in-organization",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}",
		Summary:     "Get tenant in organization",
		Description: "Get tenant by ID within an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.GetTenantInOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "create-tenant-in-organization",
		Method:      "POST",
		Path:        "/api/v1/organizations/{organizationId}/tenants",
		Summary:     "Create tenant in organization",
		Description: "Create a new tenant within an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.CreateTenantInOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "update-tenant-in-organization",
		Method:      "PUT",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}",
		Summary:     "Update tenant in organization",
		Description: "Update an existing tenant within an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.UpdateTenantInOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "delete-tenant-in-organization",
		Method:      "DELETE",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}",
		Summary:     "Delete tenant in organization",
		Description: "Delete a tenant within an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.DeleteTenantInOrganization)

	// Global tenant endpoint (for subdomain lookup during login)
	huma.Register(api, huma.Operation{
		OperationID: "get-tenant-by-subdomain",
		Method:      "GET",
		Path:        "/api/v1/tenants/subdomain/{subdomain}",
		Summary:     "Get tenant by subdomain",
		Description: "Get tenant by subdomain (used for login/redirect)",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.GetTenantBySubdomain)
}
