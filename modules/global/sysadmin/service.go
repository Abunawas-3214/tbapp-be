package sysadmin

import (
	"context"
	"errors"
	"fmt"
	"tbapp-be/common/security"
	"tbapp-be/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// DBConn didefinisikan agar Service bisa memulai transaksi (Begin)
type DBConn interface {
	db.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct {
	repo *db.Queries
	db   DBConn
}

func NewService(r *db.Queries, db DBConn) *Service {
	return &Service{repo: r, db: db}
}

func (s *Service) RegisterNewAdmin(ctx context.Context, req CreateAdminRequest) (*AdminResponse, error) {
	// 1. Hash Password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 2. Mulai Transaksi Database
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// 3. Simpan ke public.users
	user, err := qtx.CreateBaseUser(ctx, db.CreateBaseUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create base user: %w", err)
	}

	// 4. Simpan ke public.system_admins
	admin, err := qtx.CreateSystemAdmin(ctx, db.CreateSystemAdminParams{
		UserID: user.ID,
		Level:  db.AdminLevel(req.Level),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create system admin: %w", err)
	}

	// 5. Commit Transaksi
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &AdminResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Level: string(admin.Level),
	}, nil
}

func (s *Service) GetAllAdmins(ctx context.Context) ([]AdminListResponse, error) {
	// Memanggil query hasil generate SQLC
	admins, err := s.repo.ListAllSystemAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admins: %w", err)
	}

	// Mapping dari model database (db.ListAllSystemAdminsRow) ke DTO
	var response []AdminListResponse
	for _, a := range admins {
		response = append(response, AdminListResponse{
			ID:         a.ID,
			Name:       a.Name,
			Email:      a.Email,
			Level:      string(a.Level),
			IsActive:   a.IsActive,
			AdminSince: a.AdminSince.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}

func (s *Service) GetAdminByID(ctx context.Context, id string) (*AdminListResponse, error) {
	a, err := s.repo.GetSystemAdminByID(ctx, id)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	return &AdminListResponse{
		ID:         a.ID,
		Name:       a.Name,
		Email:      a.Email,
		Level:      string(a.Level),
		IsActive:   a.IsActive,
		AdminSince: a.CreatedAt.Time.Format("2006-01-02"),
	}, nil
}

func (s *Service) UpdateAdmin(ctx context.Context, id string, req UpdateAdminRequest) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)

	// 1. Ambil data lama dulu untuk backup jika data tidak dikirim
	oldData, err := qtx.GetSystemAdminByID(ctx, id)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}

	// 2. Tentukan nilai baru (Gunakan data lama jika req nil)
	finalName := oldData.Name
	if req.Name != nil {
		finalName = *req.Name
	}

	finalActive := oldData.IsActive
	if req.IsActive != nil {
		finalActive = *req.IsActive
	}

	// Jalankan Update Table Users
	err = qtx.UpdateSystemAdmin(ctx, db.UpdateSystemAdminParams{
		ID:       id,
		Name:     finalName,
		IsActive: finalActive,
	})
	if err != nil {
		return err
	}

	// 3. Update Table System Admins (jika level dikirim)
	if req.Level != nil {
		err = qtx.UpdateAdminLevel(ctx, db.UpdateAdminLevelParams{
			UserID: id,
			Level:  db.AdminLevel(*req.Level),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Service) DeleteAdmin(ctx context.Context, id string) error {
	err := s.repo.DeleteSystemAdmin(ctx, id)
	if err != nil {
		return fmt.Errorf("gatal menghapus admin: %w", err)
	}
	return nil
}

func (s *Service) ChangeAdminPassword(ctx context.Context, id string, newPassword string) error {
	// 1. Hash password baru
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 2. Update ke database
	return s.repo.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	})
}

func (s *Service) ChangeAdminEmail(ctx context.Context, id string, newEmail string) error {
	// 1. Cek apakah email baru sudah dipakai oleh user lain
	existing, _ := s.repo.GetUserByEmail(ctx, newEmail)
	if existing.ID != "" {
		// Jika ditemukan user dengan email tersebut dan ID-nya berbeda
		if existing.ID != id {
			return errors.New("email sudah digunakan oleh pengguna lain")
		}
		// Jika ID-nya sama, berarti dia mencoba mengganti ke email yang sama, biarkan sukses
		return nil
	}

	// 2. Eksekusi update email
	err := s.repo.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		ID:    id,
		Email: newEmail,
	})
	if err != nil {
		return fmt.Errorf("gagal memperbarui email: %w", err)
	}

	return nil
}
