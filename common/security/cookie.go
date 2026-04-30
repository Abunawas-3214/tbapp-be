package security

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

// SetAuthCookie menyisipkan HTTP-Only cookie ke dalam response header
func SetAuthCookie(c fiber.Ctx, tokenName string, tokenValue string) {
	cookie := new(fiber.Cookie)
	cookie.Name = tokenName
	cookie.Value = tokenValue
	cookie.Domain = os.Getenv("COOKIE_DOMAIN")
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// Keamanan ketat
	cookie.HTTPOnly = true
	cookie.Secure = os.Getenv("COOKIE_SECURE") == "true"
	cookie.SameSite = "Lax"

	// Terapkan ke context Fiber
	c.Cookie(cookie)
}

// logout
func ClearAuthCookie(c fiber.Ctx, tokenName string) {
	cookie := new(fiber.Cookie)
	cookie.Name = tokenName
	cookie.Value = ""
	cookie.Domain = os.Getenv("COOKIE_DOMAIN")
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set expired ke masa lalu agar browser menghapusnya
	cookie.HTTPOnly = true

	c.Cookie(cookie)
}
