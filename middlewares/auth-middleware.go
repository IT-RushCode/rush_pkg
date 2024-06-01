package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/utils"
)

type AuthData struct {
	UserID uint   `json:"userId"`
	IP     string `json:"ip"`
}

// AuthMiddleware представляет собой middleware для аутентификации пользователя.
type AuthMiddleware struct {
	jwtTTL time.Duration
	jwt    utils.JWTService
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Second
	jwtRTTL := time.Duration(cfg.JWT.REFRESH_TTL) * time.Second
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL)

	return &AuthMiddleware{
		jwtTTL: jwtTTL,
		jwt:    jwtService,
	}
}

// VerifyToken проверяет токен аутентификации пользователя.
func (m *AuthMiddleware) VerifyToken(ctx *fiber.Ctx) error {
	// Список маршрутов, которые не требуют проверки токена
	noAuthRoutes := []string{
		"/",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh-token",
	}

	for _, route := range noAuthRoutes {
		if ctx.Path() == route {
			return ctx.Next()
		}
	}

	// Проверка наличия токена в заголовке запроса
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorUnauthorizedResponse(ctx, "отсутствует токен авторизации", nil)
	}

	// Проверка формата токена + на наличие Bearer
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.ErrorUnauthorizedResponse(ctx, "неверный формат токена", nil)
	}

	// Извлечение самого токена
	token := parts[1]

	// Проверка токена через сервис JWT
	claims, err := m.jwt.ValidateToken(token)
	if err != nil {
		return utils.ErrorUnauthorizedResponse(ctx, "неверный токен авторизации", nil)
	}

	ctx.Locals("UserID", claims.UserID)

	return ctx.Next()
}
