package utils

import (
	"fmt"
	"math/rand"
)

// GenerateOTP генерирует 4-значный OTP-кода.
func GenerateOTP() string {
	var code string
	for i := 0; i < 4; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10)) // Генерируем случайную цифру от 0 до 9
	}
	return code
}
