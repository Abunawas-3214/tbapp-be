package tenant

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateTenant(c fiber.Ctx) error {
	req := new(CreateTenantRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	res, err := h.service.CreateTenant(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Tenant berhasil dibuat", "data": res})
}
