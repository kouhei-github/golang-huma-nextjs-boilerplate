package di

import (
	"ai-matching/src/api/auth/organization/router"
	tenantRouter "ai-matching/src/api/auth/tenant/router"
	userRouter "ai-matching/src/api/auth/user/router"
	authRouter "ai-matching/src/api/public/authentication/router"
	healthRouter "ai-matching/src/api/public/health/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRouter(container *Container) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	config := huma.DefaultConfig("Clinic RAG API", "1.0.0")
	config.DocsPath = "/docs"
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humafiber.New(app, config)

	publicAPI := app.Group("/api/v1/public")
	authAPI := app.Group("/api/v1/auth")

	healthRouter.RegisterHealthRoutes(api, publicAPI)
	authRouter.RegisterAuthRoutes(api, publicAPI, container.Queries)

	router.RegisterOrganizationRoutes(api, authAPI, container.Queries)
	tenantRouter.RegisterTenantRoutes(api, authAPI, container.Queries)
	userRouter.RegisterUserRoutes(api, authAPI, container.Queries)

	return app
}
