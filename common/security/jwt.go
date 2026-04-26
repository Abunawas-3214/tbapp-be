package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	AdminLevel string `json:"admin_level,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, name, email, level string, secretKey string, expiryHours int) (string, error) {
	claims := CustomClaims{
		UserID:     userId,
		Name:       name,
		Email:      email,
		AdminLevel: level,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
