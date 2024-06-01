package utils

import "github.com/gofiber/fiber/v2"

// GetClientIP получает реальный IP-адрес клиента, проверяя заголовки
func GetClientIP(ctx *fiber.Ctx) string {
	ip := ctx.Get("X-Forwarded-For")
	if ip == "" {
		ip = ctx.Get("X-Real-IP")
	}
	if ip == "" {
		ip = ctx.IP()
	}
	return ip
}
