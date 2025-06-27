package router

import (
	"ai-matching/src/api/auth/tenant_user/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterTenantUserRoutes(api huma.API, router fiber.Router, tenantUserController *controller.TenantUserController) {
	// Add user to tenant
	huma.Register(api, huma.Operation{
		OperationID: "add-user-to-tenant",
		Method:      "POST",
		Path:        "/api/v1/auth/tenants/{tenantId}/users",
		Summary:     "Add user to tenant",
		Description: "Add a user to a tenant with a specified role",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.AddUserToTenant)

	// Remove user from tenant
	huma.Register(api, huma.Operation{
		OperationID: "remove-user-from-tenant",
		Method:      "DELETE",
		Path:        "/api/v1/auth/tenants/{tenantId}/users/{userId}",
		Summary:     "Remove user from tenant",
		Description: "Remove a user from a tenant",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.RemoveUserFromTenant)

	// Get all users in a tenant
	huma.Register(api, huma.Operation{
		OperationID: "get-tenant-users",
		Method:      "GET",
		Path:        "/api/v1/auth/tenants/{tenantId}/users",
		Summary:     "Get tenant users",
		Description: "Get all users in a tenant",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.GetTenantUsers)

	// Get all tenants for a user
	huma.Register(api, huma.Operation{
		OperationID: "get-user-tenants",
		Method:      "GET",
		Path:        "/api/v1/auth/users/{userId}/tenants",
		Summary:     "Get user tenants",
		Description: "Get all tenants a user belongs to",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.GetUserTenants)

	// Update user role in tenant
	huma.Register(api, huma.Operation{
		OperationID: "update-user-role",
		Method:      "PUT",
		Path:        "/api/v1/auth/tenants/{tenantId}/users/{userId}/role",
		Summary:     "Update user role",
		Description: "Update a user's role in a tenant",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.UpdateUserRole)
}
