package sysadmin

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	admin := router.Group("/admins")

	admin.Post("/", h.CreateAdmin)
	admin.Get("/", h.GetListAdmin)
	admin.Get("/:id", h.GetAdmin)
	admin.Patch("/:id", h.UpdateAdmin)
	admin.Delete("/:id", h.DeleteAdmin)
	admin.Put("/:id/password", h.ChangePassword)
	admin.Put("/:id/email", h.ChangeEmail)
}
