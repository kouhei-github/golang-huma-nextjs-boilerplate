package router

import (
	"ai-matching/src/api/auth/tenant_user/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterTenantUserRoutes(api huma.API, router fiber.Router, tenantUserController *controller.TenantUserController) {
	// List users in a tenant
	huma.Register(api, huma.Operation{
		OperationID: "list-tenant-users",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}/users",
		Summary:     "List users in tenant",
		Description: "List all users in a tenant within an organization",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.ListTenantUsersInOrganization)

	// Get specific user in tenant
	huma.Register(api, huma.Operation{
		OperationID: "get-tenant-user",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}/users/{userId}",
		Summary:     "Get user in tenant",
		Description: "Get a specific user in a tenant within an organization",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.GetTenantUserInOrganization)

	// Add user to tenant
	huma.Register(api, huma.Operation{
		OperationID: "add-user-to-tenant",
		Method:      "POST",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}/users",
		Summary:     "Add user to tenant",
		Description: "Add a user to a tenant with a specified role",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.AddUserToTenantInOrganization)

	// Update user role in tenant
	huma.Register(api, huma.Operation{
		OperationID: "update-tenant-user-role",
		Method:      "PUT",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}/users/{userId}",
		Summary:     "Update user role in tenant",
		Description: "Update a user's role in a tenant",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.UpdateUserRoleInOrganization)

	// Remove user from tenant
	huma.Register(api, huma.Operation{
		OperationID: "remove-user-from-tenant",
		Method:      "DELETE",
		Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}/users/{userId}",
		Summary:     "Remove user from tenant",
		Description: "Remove a user from a tenant",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.RemoveUserFromTenantInOrganization)

	// Get all tenants for a user within an organization
	huma.Register(api, huma.Operation{
		OperationID: "get-user-tenants-in-organization",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/users/{userId}/tenants",
		Summary:     "Get user's tenants in organization",
		Description: "Get all tenants a user belongs to within an organization",
		Tags:        []string{"Tenant Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tenantUserController.GetUserTenantsInOrganization)
}
