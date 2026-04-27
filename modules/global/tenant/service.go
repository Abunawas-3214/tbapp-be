package tenant

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"tbapp-be/common/security"
	"tbapp-be/database/blueprints"
	"tbapp-be/internal/db"
	"tbapp-be/internal/tenantdb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo  *db.Queries
	tRepo *tenantdb.Queries
	pool  *pgxpool.Pool
}

func NewService(r *db.Queries, tr *tenantdb.Queries, p *pgxpool.Pool) *Service {
	return &Service{repo: r, tRepo: tr, pool: p}
}

func (s *Service) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*CreateTenantResponse, error) {
	// 1. Validasi Schema Name (Hanya huruf, angka, underscore)
	if !regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(req.SchemaName) {
		return nil, fmt.Errorf("schema name hanya boleh berisi huruf, angka, dan underscore")
	}

	// 2. Mulai Transaksi Global
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// 3. Simpan data ke skema PUBLIC
	hashedPassword, _ := security.HashPassword(req.OwnerPassword)
	owner, err := qtx.CreateBaseUser(ctx, db.CreateBaseUserParams{
		Name:         req.OwnerName,
		Email:        req.OwnerEmail,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat user owner: %w", err)
	}

	store, err := qtx.CreateStore(ctx, db.CreateStoreParams{
		Name:       req.StoreName,
		Slug:       req.StoreSlug,
		SchemaName: req.SchemaName,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat store pusat: %w", err)
	}

	err = qtx.CreateUserStoreAccess(ctx, db.CreateUserStoreAccessParams{
		UserID:  owner.ID,
		StoreID: store.ID,
	})
	if err != nil {
		return nil, err
	}

	// 4. BUAT SKEMA BARU & JALANKAN BLUEPRINT
	if _, err := tx.Exec(ctx, fmt.Sprintf("CREATE SCHEMA %s", req.SchemaName)); err != nil {
		return nil, fmt.Errorf("gagal membuat skema: %w", err)
	}

	// Arahkan pencarian ke skema baru
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL search_path TO %s", req.SchemaName)); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, blueprints.TenantBlueprintSQL); err != nil {
		return nil, fmt.Errorf("gagal migrasi tabel tenant: %w", err)
	}

	// 5. ISI DATA AWAL (ROLE & EMPLOYEE)
	// Kita gunakan tRepo yang sudah di-bind dengan transaksi saat ini
	tqtx := s.tRepo.WithTx(tx)

	roleID := uuid.New().String()
	_, err = tqtx.CreateRole(ctx, tenantdb.CreateRoleParams{
		ID:          roleID,
		Name:        "OWNER",
		Permissions: []byte(`{"all": true}`), // Contoh permission full access
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat role owner: %w", err)
	}

	_, err = tqtx.CreateEmployee(ctx, tenantdb.CreateEmployeeParams{
		ID:       uuid.New().String(),
		UserID:   owner.ID,
		FullName: req.OwnerName,
		RoleID:   roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat data employee owner: %w", err)
	}

	// 6. Selesai
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &CreateTenantResponse{
		StoreID:   store.ID,
		StoreName: store.Name,
		StoreSlug: store.Slug,
		OwnerID:   owner.ID,
		OwnerName: owner.Name,
	}, nil
}
