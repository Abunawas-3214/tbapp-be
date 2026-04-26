package sysadmin

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateAdmin(c fiber.Ctx) error {
	req := new(CreateAdminRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	res, err := h.service.RegisterNewAdmin(c.Context(), *req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Admin berhasil didaftarkan",
		"data":    res,
	})
}

func (h *Handler) GetListAdmin(c fiber.Ctx) error {
	admins, err := h.service.GetAllAdmins(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil daftar admin"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Daftar admin berhasil diambil",
		"data":    admins,
	})
}

func (h *Handler) GetAdmin(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetAdminByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *Handler) UpdateAdmin(c fiber.Ctx) error {
	id := c.Params("id")
	req := new(UpdateAdminRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateAdmin(c.Context(), id, *req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Admin berhasil diperbarui"})
}

func (h *Handler) DeleteAdmin(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteAdmin(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Admin berhasil dihapus"})
}

func (h *Handler) ChangePassword(c fiber.Ctx) error {
	id := c.Params("id")
	req := new(ChangePasswordRequest)

	// 1. Parsing Body
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// 2. Validasi kecocokan password (Server-side validation)
	if req.NewPassword != req.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Password baru dan konfirmasi tidak cocok"})
	}

	// 3. Validasi panjang password minimal (Contoh: 8 karakter)
	if len(req.NewPassword) < 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Password minimal harus 8 karakter"})
	}

	// 4. Panggil Service
	if err := h.service.ChangeAdminPassword(c.Context(), id, req.NewPassword); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengganti password: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Password admin berhasil diperbarui",
	})
}

func (h *Handler) ChangeEmail(c fiber.Ctx) error {
	id := c.Params("id")
	req := new(ChangeEmailRequest)

	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	if err := h.service.ChangeAdminEmail(c.Context(), id, req.NewEmail); err != nil {
		// Kita berikan status 409 Conflict jika email sudah ada
		return c.Status(409).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Email admin berhasil diperbarui",
	})
}
