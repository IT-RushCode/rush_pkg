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
	jwtTTL          time.Duration
	jwt             utils.JWTService
	whiteListRoutes []string
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(cfg *config.Config, whitelistRoutes []string) *AuthMiddleware {
	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Second
	jwtRTTL := time.Duration(cfg.JWT.REFRESH_TTL) * time.Second
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL)

	return &AuthMiddleware{
		jwtTTL:          jwtTTL,
		jwt:             jwtService,
		whiteListRoutes: whitelistRoutes,
	}
}

// VerifyToken проверяет токен аутентификации пользователя.
func (m *AuthMiddleware) VerifyToken(ctx *fiber.Ctx) error {

	for _, route := range m.whiteListRoutes {
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

	// Проверка привилегий пользователя
	// if err := m.checkPermissions(ctx, claims.UserID); err != nil {
	// 	return err
	// }

	ctx.Locals("UserID", claims.UserID)

	return ctx.Next()
}
