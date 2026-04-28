package store_profile

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	profile := router.Group("profile")
	profile.Get("/", h.GetProfile)

}
