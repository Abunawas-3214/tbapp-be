package main

import (
	"log"
	"os"
	"tbapp-be/config"
	"tbapp-be/internal/db"
	"tbapp-be/internal/tenantdb"
	"tbapp-be/middleware"

	_ "tbapp-be/docs"

	"tbapp-be/modules/app/store_profile"
	"tbapp-be/modules/global/auth"
	"tbapp-be/modules/global/sysadmin"
	"tbapp-be/modules/global/tenant"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

// @title TB App API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/v1
func main() {
	// Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 1. Connect ke Neon
	dbPool := config.ConnectDB()
	defer dbPool.Close()

	// 2. Inisialisasi SQLC Queries
	repoPublic := db.New(dbPool)
	repoTenant := tenantdb.New(dbPool)

	// 3. Setup Fiber
	app := fiber.New()
	api := app.Group("/api/v1")
	globalGroup := api.Group("/")
	adminGroup := api.Group("/admin")
	tenantGroup := api.Group("/app/:store_slug")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.RoleAdminGuard())
	tenantGroup.Use(middleware.TenantMiddleware(dbPool))

	// Module: Auth
	authService := auth.NewService(repoPublic, repoTenant, dbPool)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(globalGroup)

	// Module: Sysadmin
	adminService := sysadmin.NewService(repoPublic, dbPool)
	adminHandler := sysadmin.NewHandler(adminService)
	adminHandler.RegisterRoutes(adminGroup)

	// Module: Tenant
	tenantService := tenant.NewService(repoPublic, repoTenant, dbPool)
	tenantHandler := tenant.NewHandler(tenantService)
	tenantHandler.RegisterRoutes(adminGroup)

	// Module: Store Profile
	storeProfileService := store_profile.NewService(repoPublic, repoTenant, dbPool)
	storeProfileHandler := store_profile.NewHandler(storeProfileService)
	storeProfileHandler.RegisterRoutes(tenantGroup)

	// Swagger
	cfg := swagger.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "TB App API Documentation",
	}

	app.Use(swagger.New(cfg))

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
