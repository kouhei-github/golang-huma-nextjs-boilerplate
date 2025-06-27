package router

import (
	"ai-matching/src/api/auth/tenant/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterTenantRoutes(api huma.API, router fiber.Router, tenantController *controller.TenantController) {

	huma.Register(api, huma.Operation{
		OperationID: "get-tenant",
		Method:      "GET",
		Path:        "/api/v1/auth/tenants/{id}",
		Summary:     "Get tenant",
		Description: "Get tenant by ID",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.GetTenant)

	huma.Register(api, huma.Operation{
		OperationID: "get-tenant-by-subdomain",
		Method:      "GET",
		Path:        "/api/v1/auth/tenants/subdomain/{subdomain}",
		Summary:     "Get tenant by subdomain",
		Description: "Get tenant by subdomain",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.GetTenantBySubdomain)

	huma.Register(api, huma.Operation{
		OperationID: "list-tenants-by-organization",
		Method:      "GET",
		Path:        "/api/v1/auth/organizations/{organizationId}/tenants",
		Summary:     "List tenants by organization",
		Description: "List all tenants for an organization",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.ListTenantsByOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "create-tenant",
		Method:      "POST",
		Path:        "/api/v1/auth/tenants",
		Summary:     "Create tenant",
		Description: "Create a new tenant",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.CreateTenant)

	huma.Register(api, huma.Operation{
		OperationID: "update-tenant",
		Method:      "PUT",
		Path:        "/api/v1/auth/tenants/{id}",
		Summary:     "Update tenant",
		Description: "Update an existing tenant",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.UpdateTenant)

	huma.Register(api, huma.Operation{
		OperationID: "delete-tenant",
		Method:      "DELETE",
		Path:        "/api/v1/auth/tenants/{id}",
		Summary:     "Delete tenant",
		Description: "Delete a tenant",
		Tags:        []string{"Tenants"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantController.DeleteTenant)
}
