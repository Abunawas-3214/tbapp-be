package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

func SetupSecurity(app *fiber.App) {
	app.Use(helmet.New())

	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		originsEnv = "http://localhost:3000"
	}

	allowedOrigins := strings.Split(originsEnv, ",")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		MaxAge: 3600,
	}))
}
