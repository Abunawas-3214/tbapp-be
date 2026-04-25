package security

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah plain text menjadi hash bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12) // cost 12 sudah cukup aman
	return string(bytes), err
}

// CheckPasswordHash membandingkan plain text dengan hash (untuk login nanti)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
