package middleware

import (
	"fmt"
	"os"
	"strings"
	"tbapp-be/common/security"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TenantMiddleware(dbPool *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Ambil slug dari URL Param (:store_slug)
		urlSlug := c.Params("store_slug")
		if urlSlug == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Slug toko diperlukan"})
		}

		// 2. Ambil dan Verifikasi Token
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenString, &security.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Sesi toko tidak valid"})
		}

		claims, _ := token.Claims.(*security.CustomClaims)

		// 3. KEAMANAN KRUSIAL: Cross-Validation
		// Cek apakah slug di URL sama dengan slug yang diizinkan di Token
		if claims.StoreSlug == nil || *claims.StoreSlug != urlSlug {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses dilarang: Token tidak valid untuk toko ini",
			})
		}

		// 4. Switch Database Path
		query := fmt.Sprintf("SET search_path TO %s, public", *claims.SchemaName)
		if _, err := dbPool.Exec(c.Context(), query); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengalihkan database"})
		}

		// Simpan data untuk kebutuhan audit log di level controller
		c.Locals("user_id", claims.UserID)
		c.Locals("store_id", *claims.StoreID)

		return c.Next()
	}
}
