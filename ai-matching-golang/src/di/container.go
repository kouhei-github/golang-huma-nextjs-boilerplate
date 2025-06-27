package di

import (
	db "ai-matching/db/sqlc"
	authController "ai-matching/src/api/auth/organization/controller"
	organizationUsecase "ai-matching/src/api/auth/organization/usecase"
	tenantController "ai-matching/src/api/auth/tenant/controller"
	tenantUsecase "ai-matching/src/api/auth/tenant/usecase"
	tenantUserController "ai-matching/src/api/auth/tenant_user/controller"
	tenantUserUsecase "ai-matching/src/api/auth/tenant_user/usecase"
	userController "ai-matching/src/api/auth/user/controller"
	userUsecase "ai-matching/src/api/auth/user/usecase"
	publicAuthController "ai-matching/src/api/public/authentication/controller"
	publicAuthUsecase "ai-matching/src/api/public/authentication/usecase"
	healthController "ai-matching/src/api/public/health/controller"
	"ai-matching/src/domain/interface/external"
	"ai-matching/src/domain/interface/repository"
	"ai-matching/src/infrastructure/external/cognito"
	infraRepository "ai-matching/src/infrastructure/repository"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Container struct {
	DB            *sqlx.DB
	Queries       db.Querier
	CognitoClient external.CognitoClient

	// Repositories
	UserRepository         repository.UserRepository
	OrganizationRepository repository.OrganizationRepository
	TenantRepository       repository.TenantRepository
	TenantUserRepository   repository.TenantUserRepository

	// Usecases
	AuthUsecase         *publicAuthUsecase.AuthUsecase
	UserUsecase         *userUsecase.UserUsecase
	OrganizationUsecase *organizationUsecase.OrganizationUsecase
	TenantUsecase       *tenantUsecase.TenantUsecase
	TenantUserUsecase   *tenantUserUsecase.TenantUserUsecase

	// Controllers
	AuthController         *publicAuthController.AuthController
	UserController         *userController.UserController
	OrganizationController *authController.OrganizationController
	TenantController       *tenantController.TenantController
	TenantUserController   *tenantUserController.TenantUserController
	HealthController       *healthController.HealthController
}

func NewContainer() *Container {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		sslMode := os.Getenv("DB_SSL_MODE")

		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		if sslMode == "" {
			sslMode = "disable"
		}

		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	queries := db.New(sqlDB)

	cognitoClient, err := cognito.NewCognitoClient()
	if err != nil {
		log.Fatal("Failed to create Cognito client:", err)
	}

	// Initialize repositories
	userRepo := infraRepository.NewUserRepository(queries)
	orgRepo := infraRepository.NewOrganizationRepository(queries)
	tenantRepo := infraRepository.NewTenantRepository(queries)
	tenantUserRepo := infraRepository.NewTenantUserRepository(queries)

	// Initialize usecases
	authUc := publicAuthUsecase.NewAuthUsecase(userRepo, tenantUserRepo, tenantRepo, orgRepo, cognitoClient)
	userUc := userUsecase.NewUserUsecase(userRepo, tenantUserRepo, cognitoClient)
	orgUc := organizationUsecase.NewOrganizationUsecase(orgRepo)
	tenantUc := tenantUsecase.NewTenantUsecase(tenantRepo)
	tenantUserUc := tenantUserUsecase.NewTenantUserUsecase(tenantUserRepo, tenantRepo, userRepo)

	// Initialize controllers
	authCtrl := publicAuthController.NewAuthController(authUc)
	userCtrl := userController.NewUserController(userUc)
	orgCtrl := authController.NewOrganizationController(orgUc)
	tenantCtrl := tenantController.NewTenantController(tenantUc)
	tenantUserCtrl := tenantUserController.NewTenantUserController(tenantUserUc)
	healthCtrl := healthController.NewHealthController()

	return &Container{
		DB:            sqlxDB,
		Queries:       queries,
		CognitoClient: cognitoClient,

		// Repositories
		UserRepository:         userRepo,
		OrganizationRepository: orgRepo,
		TenantRepository:       tenantRepo,
		TenantUserRepository:   tenantUserRepo,

		// Usecases
		AuthUsecase:         authUc,
		UserUsecase:         userUc,
		OrganizationUsecase: orgUc,
		TenantUsecase:       tenantUc,
		TenantUserUsecase:   tenantUserUc,

		// Controllers
		AuthController:         authCtrl,
		UserController:         userCtrl,
		OrganizationController: orgCtrl,
		TenantController:       tenantCtrl,
		TenantUserController:   tenantUserCtrl,
		HealthController:       healthCtrl,
	}
}
