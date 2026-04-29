package middleware

import "github.com/gofiber/fiber/v3"

func RoleAdminGuard() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Ambil admin_level yang sudah diset oleh AuthMiddleware sebelumnya
		adminLevel := c.Locals("admin_level")

		// Jika adminLevel nil atau string kosong, berarti dia user biasa/toko
		if adminLevel == nil || adminLevel == "" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses Ditolak: Area ini hanya untuk Administrator Sistem",
			})
		}

		return c.Next()
	}
}
