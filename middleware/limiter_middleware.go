package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// GlobalLimiter membatasi traffic umum agar server tetap stabil
func GlobalLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // Maksimal 100 request
		Expiration: 1 * time.Minute, // Per 1 menit
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP() // Identifikasi berdasarkan IP
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Terlalu banyak permintaan, silakan coba lagi nanti.",
			})
		},
	})
}

// AuthLimiter jauh lebih ketat untuk mencegah Brute Force pada Login
func AuthLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,                // Hanya 5 percobaan
		Expiration: 15 * time.Minute, // Jika habis, tunggu 15 menit
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Terlalu banyak percobaan login. Akun Anda ditangguhkan selama 15 menit.",
			})
		},
	})
}
