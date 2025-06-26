package router

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/organization/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrganizationRoutes(api huma.API, router fiber.Router, queries db.Querier) {
	orgController := controller.NewOrganizationController(queries)

	huma.Register(api, huma.Operation{
		OperationID: "get-organization",
		Method:      "GET",
		Path:        "/api/v1/auth/organizations/{id}",
		Summary:     "Get organization",
		Description: "Get organization by ID",
		Tags:        []string{"Organizations"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, orgController.GetOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "list-organizations",
		Method:      "GET",
		Path:        "/api/v1/auth/organizations",
		Summary:     "List organizations",
		Description: "List all organizations",
		Tags:        []string{"Organizations"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, orgController.ListOrganizations)

	huma.Register(api, huma.Operation{
		OperationID: "create-organization",
		Method:      "POST",
		Path:        "/api/v1/auth/organizations",
		Summary:     "Create organization",
		Description: "Create a new organization",
		Tags:        []string{"Organizations"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, orgController.CreateOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "update-organization",
		Method:      "PUT",
		Path:        "/api/v1/auth/organizations/{id}",
		Summary:     "Update organization",
		Description: "Update an existing organization",
		Tags:        []string{"Organizations"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, orgController.UpdateOrganization)

	huma.Register(api, huma.Operation{
		OperationID: "delete-organization",
		Method:      "DELETE",
		Path:        "/api/v1/auth/organizations/{id}",
		Summary:     "Delete organization",
		Description: "Delete an organization",
		Tags:        []string{"Organizations"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, orgController.DeleteOrganization)
}
