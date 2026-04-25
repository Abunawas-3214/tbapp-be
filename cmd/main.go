package main

import (
	"log"
	"tbapp-be/config"
	"tbapp-be/internal"
	"tbapp-be/modules/users"

	_ "tbapp-be/docs"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v3"
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
	// 1. Connect ke Neon
	dbPool := config.ConnectDB()
	defer dbPool.Close()

	// 2. Inisialisasi SQLC Queries
	queries := internal.New(dbPool)

	// 3. Setup Fiber
	app := fiber.New()
	api := app.Group("/api/v1")

	// 4. Daftarkan Modul Users
	userHandler := &users.UserHandler{Queries: queries}
	users.InitUserRoutes(api, userHandler)

	// Swagger
	cfg := swagger.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "TB App API Documentationi",
	}

	app.Use(swagger.New(cfg))

	// 5. Jalankan Server
	log.Fatal(app.Listen(":3000"))
}
