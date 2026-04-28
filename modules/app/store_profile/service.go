package store_profile

import (
	"context"
	"tbapp-be/internal/db"
	"tbapp-be/internal/tenantdb"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pRepo  *db.Queries
	tRepo  *tenantdb.Queries
	dbPool *pgxpool.Pool
}

func NewService(pr *db.Queries, tr *tenantdb.Queries, db *pgxpool.Pool) *Service {
	return &Service{pRepo: pr, tRepo: tr, dbPool: db}
}

func toText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func (s *Service) GetProfile(ctx context.Context, storeID string) (*StoreProfileResponse, error) {
	profile, err := s.tRepo.GetStoreProfile(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return &StoreProfileResponse{
		ID:          profile.StoreID,
		Name:        profile.StoreName,
		Description: &profile.Description.String,
		Address:     &profile.Address.String,
		Phone:       &profile.Phone.String,
		Email:       &profile.Email.String,
		TaxID:       &profile.TaxID.String,
		LogoURL:     &profile.LogoUrl.String,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, req UpdateProfileRequest, storeID string) (bool, error) {
	// 1. Mulai Transaksi
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return false, err
	}
	// Pastikan rollback jika terjadi error sebelum commit
	defer tx.Rollback(ctx)

	// 2. Suntikkan transaksi ke repository
	qPublic := s.pRepo.WithTx(tx)
	qTenant := s.tRepo.WithTx(tx)

	// 3. Update skema PUBLIC (Nama Toko)
	// Kita berasumsi Name dikirim (baik berubah atau tidak)
	err = qPublic.UpdateStoreName(ctx, db.UpdateStoreNameParams{
		ID:   storeID,
		Name: req.Name,
	})
	if err != nil {
		return false, err
	}

	// 4. Update skema TENANT (Detail Profil)
	// Menggunakan helper toText untuk mencegah panic dereference
	_, err = qTenant.UpdateStoreProfile(ctx, tenantdb.UpdateStoreProfileParams{
		Description: toText(req.Description),
		Address:     toText(req.Address),
		Phone:       toText(req.Phone),
		Email:       toText(req.Email),
		TaxID:       toText(req.TaxID),
		LogoUrl:     toText(req.LogoURL),
		StoreID:     storeID,
	})
	if err != nil {
		return false, err
	}

	// 5. Eksekusi COMMIT
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
