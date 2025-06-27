package router

import (
	"ai-matching/src/api/auth/user/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(api huma.API, router fiber.Router, userController *controller.UserController) {

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      "GET",
		Path:        "/api/v1/auth/users/{userId}",
		Summary:     "Get user",
		Description: "Get user by ID",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.GetUser)

	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      "GET",
		Path:        "/api/v1/auth/users",
		Summary:     "List users",
		Description: "List all users",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      "POST",
		Path:        "/api/v1/auth/users",
		Summary:     "Create user",
		Description: "Create a new user",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      "PUT",
		Path:        "/api/v1/auth/users/{userId}",
		Summary:     "Update user",
		Description: "Update an existing user",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      "DELETE",
		Path:        "/api/v1/auth/users/{userId}",
		Summary:     "Delete user",
		Description: "Delete a user",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, userController.DeleteUser)
}
