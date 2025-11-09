package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

// CalculateFileHash считает SHA256 хэш файла
func CalculateFileHash(file multipart.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // вернём указатель в начало
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
