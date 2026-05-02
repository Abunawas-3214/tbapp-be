package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"tbapp-be/common/security"
	"tbapp-be/internal/db"

	"tbapp-be/internal/tenantdb"

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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai sesi validasi toko: %w", err)
	}
	defer tx.Rollback(ctx)

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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Set search_path ke skema toko agar query sqlc mencari ke tabel yang benar
	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL search_path TO %q", selectedStore.SchemaName))
	if err != nil {
		return nil, fmt.Errorf("gagal berpindah ke skema toko: %w", err)
	}

	// Gunakan repository dengan transaksi (WithTx)
	tqtx := s.tRepo.WithTx(tx)

	// Ambil data Role & Permission dari tabel employees & roles milik tenant
	empRole, err := tqtx.GetTenantEmployeeRole(ctx, userID)
	if err != nil {
		// Jika tidak ditemukan di tabel employees tenant tersebut
		return nil, fmt.Errorf("Gagal mengambil data karyawan (DB Error: %w)", err)
	}

	// Parsing JSON permissions dari database ([]byte) ke map Go
	var permissions map[string]interface{}
	if err := json.Unmarshal(empRole.Permissions, &permissions); err != nil {
		permissions = make(map[string]interface{})
	}

	// 3. Generate Tenant Token
	secretKey := os.Getenv("JWT_SECRET")
	token, err := security.GenerateTenantToken(
		user.ID,
		user.Name,
		user.Email,
		empRole.RoleName,
		permissions,
		selectedStore.ID,
		selectedStore.Slug,
		selectedStore.SchemaName,
		secretKey,
		24,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan token toko: %w", err)
	}

	// Commit transaksi (meskipun hanya SELECT, praktik ini menjaga koneksi tetap bersih)
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &SelectStoreResponse{
		Token:      token,
		StoreName:  selectedStore.Name,
		SchemaName: selectedStore.SchemaName,
		Message:    fmt.Sprintf("Berhasil masuk ke ruang kerja %s", selectedStore.Name),
	}, nil
}
