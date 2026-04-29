package middleware

import (
	"os"
	"strings"
	"tbapp-be/common/security"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Ambil Header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}

		// 2. Format harus: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{"error": "Format token salah"})
		}

		// 3. Validasi & Parse Token
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenString, &security.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid atau sudah kadaluwarsa"})
		}

		// 4. Ambil Claims dan titipkan UserID ke Locals
		claims, ok := token.Claims.(*security.CustomClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Gagal membaca data token"})
		}

		// INI YANG DICARI HANDLER TADI:
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_name", claims.Name)
		c.Locals("admin_level", claims.AdminLevel)

		return c.Next()
	}
}
