package router

import (
	"ai-matching/src/api/auth/user/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(api huma.API, router fiber.Router, userController *controller.UserController) {
	// Organization-level user endpoints
	huma.Register(api, huma.Operation{
		OperationID: "list-organization-users",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/users",
		Summary:     "List users in organization",
		Description: "List all users in an organization",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.ListOrganizationUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-organization-user",
		Method:      "GET",
		Path:        "/api/v1/organizations/{organizationId}/users/{userId}",
		Summary:     "Get user in organization",
		Description: "Get user by ID within an organization",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.GetOrganizationUser)

	huma.Register(api, huma.Operation{
		OperationID: "create-organization-user",
		Method:      "POST",
		Path:        "/api/v1/organizations/{organizationId}/users",
		Summary:     "Create user in organization",
		Description: "Create a new user within an organization",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.CreateOrganizationUser)

	huma.Register(api, huma.Operation{
		OperationID: "update-organization-user",
		Method:      "PUT",
		Path:        "/api/v1/organizations/{organizationId}/users/{userId}",
		Summary:     "Update user in organization",
		Description: "Update an existing user within an organization",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.UpdateOrganizationUser)

	huma.Register(api, huma.Operation{
		OperationID: "delete-organization-user",
		Method:      "DELETE",
		Path:        "/api/v1/organizations/{organizationId}/users/{userId}",
		Summary:     "Delete user from organization",
		Description: "Delete a user from an organization",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.DeleteOrganizationUser)
}
