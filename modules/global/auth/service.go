package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"tbapp-be/common/security"
	"tbapp-be/internal/db"
)

type Service struct {
	repo *db.Queries
}

func NewService(r *db.Queries) *Service {
	return &Service{repo: r}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Kredensial yang Anda masukkan salah")
	}

	if !user.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan, hubungi superadmin")
	}

	match := security.CheckPasswordHash(req.Password, user.PasswordHash.String)
	if !match {
		return nil, errors.New("kredensial yang Anda masukkan salah")
	}

	var levelPtr *string
	adminAccess, _ := s.repo.GetAdminAccess(ctx, user.ID)
	if adminAccess != "" {
		levelStr := string(adminAccess)
		levelPtr = &levelStr
	}

	stores, _ := s.repo.GetUserStores(ctx, user.ID)
	var storeDTOs []StoreAccessDTO
	for _, st := range stores {
		storeDTOs = append(storeDTOs, StoreAccessDTO{
			ID:         st.ID,
			Name:       st.Name,
			Slug:       st.Slug,
			SchemaName: st.SchemaName,
		})
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET belum diset di .env")
	}

	token, err := security.GenerateToken(
		user.ID,
		user.Name,
		user.Email,
		levelPtr,
		secretKey,
		24,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan akses: %w", err)
	}

	res := &LoginResponse{
		Token: token,
	}
	res.User.ID = user.ID
	res.User.Name = user.Name
	res.User.Email = user.Email
	res.User.AdminLevel = levelPtr
	res.Stores = storeDTOs

	return res, nil
}

func (s *Service) SelectStore(ctx context.Context, userID string, req SelectStoreRequest) (*SelectStoreResponse, error) {
	// 1. Dapatkan data user (untuk info name & email di token)
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// 2. VALIDASI KEAMANAN: Cek apakah user punya akses ke toko ini
	// Kita akan loop stores milik user (bisa juga pakai query DB khusus agar lebih cepat)
	stores, err := s.repo.GetUserStores(ctx, userID)
	if err != nil {
		return nil, errors.New("gagal memuat akses toko")
	}

	var selectedStore *db.GetUserStoresRow // Asumsi nama struct Anda adalah db.Store dari hasil sqlc
	for _, st := range stores {
		if st.ID == req.StoreID {
			selectedStore = &st
			break
		}
	}

	if selectedStore == nil {
		// Logika Anti-Bypass
		return nil, errors.New("anda tidak memiliki izin untuk mengakses toko ini")
	}

	// 3. Generate Tenant Token
	secretKey := os.Getenv("JWT_SECRET")
	token, err := security.GenerateTenantToken(
		user.ID,
		user.Name,
		user.Email,
		selectedStore.ID,
		selectedStore.SchemaName,
		secretKey,
		24,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan token toko: %w", err)
	}

	return &SelectStoreResponse{
		Token:      token,
		SchemaName: selectedStore.SchemaName,
		Message:    fmt.Sprintf("Berhasil masuk ke ruang kerja %s", selectedStore.Name),
	}, nil
}
