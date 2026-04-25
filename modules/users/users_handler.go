package users

import (
	"context"
	"tbapp-be/internal"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	Queries *internal.Queries
}

// CreateUser handles: POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// 1. DTO: Penampung input
	type CreateRequest struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		RoleID string `json:"role_id"`
	}

	req := new(CreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// 2. Simpan ke database menggunakan hasil generate SQLC
	user, err := h.Queries.CreateUserByAdmin(context.Background(), internal.CreateUserByAdminParams{
		ID:     uuid.New().String(),
		Name:   req.Name,
		Email:  req.Email,
		RoleID: pgtype.Text{String: req.RoleID, Valid: req.RoleID != ""},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan user",
			"error":   err.Error(),
		})
	}

	// 3. Kirim balik data user yang baru dibuat
	return c.Status(201).JSON(user)
}

// ListUser handles: GET /api/v1/users
func (h *UserHandler) ListUser(c *fiber.Ctx) error {
	users, err := h.Queries.ListUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}

	return c.JSON(users)
}

// GetUser handles: GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id") // Ambil ID dari URL

	user, err := h.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(user)
}

// UpdateUser handles: PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateRequest struct {
		Name     string `json:"name"`
		RoleID   string `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}

	req := new(UpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	user, err := h.Queries.UpdateUser(context.Background(), internal.UpdateUserParams{
		ID:       id,
		Name:     req.Name,
		RoleID:   pgtype.Text{String: req.RoleID, Valid: req.RoleID != ""},
		IsActive: pgtype.Bool{Bool: req.IsActive, Valid: true},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// DeleteUser handles: DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.Queries.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}
