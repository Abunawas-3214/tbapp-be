package users

import "github.com/gofiber/fiber/v3"

func InitUserRoutes(router fiber.Router, h *UserHandler) {
	// Endpoint: /api/v1/users
	u := router.Group("/users")

	u.Post("/", h.CreateUser)
	u.Get("/", h.ListUser)
	u.Get("/:id", h.GetUser)
	u.Put("/:id", h.UpdateUser)
	u.Delete("/:id", h.DeleteUser)

}
