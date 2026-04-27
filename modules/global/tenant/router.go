package tenant

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	tenant := router.Group("tenants")

	tenant.Post("/", h.CreateTenant)
}
