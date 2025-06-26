package router

import (
	"ai-matching/src/api/public/health/controller"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(api huma.API, router fiber.Router) {
	healthController := controller.NewHealthController()

	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      "GET",
		Path:        "/api/v1/public/health",
		Summary:     "Health check",
		Description: "Check if the service is healthy",
		Tags:        []string{"Health"},
	}, healthController.GetHealth)
}
