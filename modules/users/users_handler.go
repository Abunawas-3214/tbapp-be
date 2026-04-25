package users

import (
	"context"
	"tbapp-be/common/security"
	"tbapp-be/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	Queries *internal.Queries
}

// DTO untuk Request
type CreateUserRequest struct {
	Name     string  `json:"name" example:"Abunawas"`
	Email    string  `json:"email" example:"abu@example.com"`
	RoleID   string  `json:"role_id" example:"ADMIN"`
	Password *string `json:"password,omitempty" example:"rahasia123"`
}

// CreateUser godoc
// @Summary      Mendaftarkan user baru (oleh Admin)
// @Description  Admin mengundang user baru. Jika password diisi, akan di-hash.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "Data User"
// @Success      201 {object} internal.User
// @Router       /users [post]
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// Logika Hashing Password
	var hashedPw pgtype.Text
	if req.Password != nil && *req.Password != "" {
		hash, err := security.HashPassword(*req.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal mengamankan password"})
		}
		hashedPw = pgtype.Text{String: hash, Valid: true}
	}

	// Eksekusi SQLC
	user, err := h.Queries.CreateUserByAdmin(context.Background(), internal.CreateUserByAdminParams{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		RoleID:       pgtype.Text{String: req.RoleID, Valid: req.RoleID != ""},
		PasswordHash: hashedPw,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(user)
}

// ListUser handles: GET /api/v1/users
func (h *UserHandler) ListUser(c fiber.Ctx) error {
	users, err := h.Queries.ListUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}

	return c.JSON(users)
}

// GetUser handles: GET /api/v1/users/:id
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id") // Ambil ID dari URL

	user, err := h.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(user)
}

// UpdateUser handles: PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")

	type UpdateRequest struct {
		Name     string `json:"name"`
		RoleID   string `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}

	req := new(UpdateRequest)
	if err := c.Bind().Body(req); err != nil {
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
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	err := h.Queries.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}
