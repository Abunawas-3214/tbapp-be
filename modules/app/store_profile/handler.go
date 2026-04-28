package store_profile

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetProfile(c fiber.Ctx) error {
	storeID, ok := c.Locals("store_id").(string)
	if !ok || storeID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Identitas toko tidak dalam sesi"})
	}

	res, err := h.service.GetProfile(c.Context(), storeID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profil toko tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	storeID := c.Locals("store_id").(string)

	var req UpdateProfileRequest // Buat struct ini di dto.go
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	_, err := h.service.UpdateProfile(c.Context(), req, storeID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui profil"})
	}

	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui"})
}
