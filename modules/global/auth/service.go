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

	adminAccess, _ := s.repo.GetAdminAccess(ctx, user.ID)
	var levelStr string
	if adminAccess != "" {
		levelStr = string(adminAccess)
	}
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET belum diset di .env")
	}

	token, err := security.GenerateToken(
		user.ID,
		user.Name,
		user.Email,
		levelStr,
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

	return res, nil
}
