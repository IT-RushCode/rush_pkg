package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt генерирует случайную соль
func GenerateSalt() string {
	salt := make([]byte, 16)
	rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt)
}

// HashPasswordWithSalt хеширует пароль с солью
func HashPasswordWithSalt(password, salt string) (string, error) {
	saltedPassword := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// ComparePasswordWithSalt сравнивает хешированный пароль с солью
func ComparePasswordWithSalt(hashedPassword, password, salt string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}
