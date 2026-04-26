package auth

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)
}
