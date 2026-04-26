package auth

import (
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

// NewHandler inisialisasi handler dengan service yang sudah di-inject
func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// Login menangani POST request untuk autentikasi
func (h *Handler) Login(c fiber.Ctx) error {
	// 1. Inisialisasi struct DTO untuk menampung request body
	req := new(LoginRequest)

	// 2. Parsing JSON ke struct
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Format data yang dikirim tidak valid",
		})
	}

	// 3. Panggil Service Login
	// Kita teruskan Context dari Fiber agar bisa digunakan oleh pgx
	res, err := h.service.Login(c.Context(), *req)
	if err != nil {
		// Jika login gagal (email salah, password salah, atau user tidak aktif)
		// Kita berikan status 401 Unauthorized
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Jika berhasil, kirim response sukses beserta Token JWT
	return c.Status(200).JSON(fiber.Map{
		"message": "Login berhasil, selamat datang!",
		"data":    res,
	})
}
