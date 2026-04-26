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

type UpdateUserRequest struct {
	Name     string  `json:"name" example:"Abunawas Updated"`
	Email    string  `json:"email" example:"abu_new@example.com"`
	RoleID   string  `json:"role_id" example:"USER"`
	Password *string `json:"password,omitempty" example:"newsecret123"`
	IsActive bool    `json:"is_active" example:"true"`
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

// ListUser godoc
// @Summary      Mengambil daftar semua user
// @Description  Menampilkan semua user yang terdaftar di sistem.
// @Tags         users
// @Produce      json
// @Success      200 {array} internal.User
// @Router       /users [get]
func (h *UserHandler) ListUser(c fiber.Ctx) error {
	users, err := h.Queries.ListUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}

	return c.JSON(users)
}

// GetUser godoc
// @Summary      Mengambil detail user berdasarkan ID
// @Description  Menampilkan informasi lengkap user tertentu.
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} internal.User
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id") // Ambil ID dari URL

	user, err := h.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(user)
}

// UpdateUser godoc
// @Summary      Memperbarui data user
// @Description  Mengubah informasi nama, role, atau status aktif user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        request body UpdateUserRequest true "Data Update"
// @Success      200 {object} internal.User
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")

	req := new(UpdateUserRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// 1. Ambil data user yang sudah ada
	existingUser, err := h.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// 2. Logika Hashing Password (hanya jika password baru diisi)
	hashedPw := existingUser.PasswordHash
	if req.Password != nil && *req.Password != "" {
		hash, err := security.HashPassword(*req.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal mengamankan password baru"})
		}
		hashedPw = pgtype.Text{String: hash, Valid: true}
	}

	// 3. Update data
	user, err := h.Queries.UpdateUser(context.Background(), internal.UpdateUserParams{
		ID:           id,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPw,
		RoleID:       pgtype.Text{String: req.RoleID, Valid: req.RoleID != ""},
		IsActive:     pgtype.Bool{Bool: req.IsActive, Valid: true},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// DeleteUser godoc
// @Summary      Menghapus user
// @Description  Menghapus user dari sistem berdasarkan ID.
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	err := h.Queries.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}
