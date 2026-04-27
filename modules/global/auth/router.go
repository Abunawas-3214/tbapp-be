package auth

import (
	"tbapp-be/middleware"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)
	auth.Post("/select-store", middleware.AuthMiddleware(), h.SelectStore)
}
