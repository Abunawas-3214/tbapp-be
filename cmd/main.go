package main

import (
	"log"
	"tbapp-be/config"
	"tbapp-be/internal"
	"tbapp-be/modules/users" // Jangan lupa import modulnya

	"github.com/gofiber/fiber/v2"
)

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

	// 5. Jalankan Server
	log.Fatal(app.Listen(":3000"))
}
