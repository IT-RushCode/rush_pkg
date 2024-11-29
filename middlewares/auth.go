package middlewares

import (
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware представляет собой middleware для аутентификации пользователя.
type AuthMiddleware struct {
	jwt          utils.JWTService
	publicRoutes map[string][]string
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(
	cfg *config.Config,
	routes map[string][]string,
) *AuthMiddleware {
	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Second
	jwtRTTL := time.Duration(cfg.JWT.REFRESH_TTL) * time.Second
	return &AuthMiddleware{
		jwt:          utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL),
		publicRoutes: routes,
	}
}

// VerifyToken выполняет основную проверку токена.
func (m *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	// Проверка на публичный маршрут
	if m.isPublicRoute(ctx) {
		ctx.Locals("IsPublic", true)

		// Если заголовок Authorization есть, пытаемся извлечь токен
		claims, _ := m.extractTokenClaims(ctx)
		if claims != nil {
			ctx.Locals("UserID", claims.UserID)
			ctx.Locals("IsMob", claims.IsMob)
		}

		return ctx.Next()
	}

	// Обычная авторизация для приватных маршрутов
	claims, err := m.extractTokenClaims(ctx)
	if err != nil {
		return err
	}

	// Сохраняем данные пользователя в контексте
	ctx.Locals("UserID", claims.UserID)
	ctx.Locals("IsMob", claims.IsMob)

	return ctx.Next()
}

// isPublicRoute проверяет, является ли маршрут и метод запросом в белом списке.
func (m *AuthMiddleware) isPublicRoute(ctx *fiber.Ctx) bool {
	path := ctx.Path()

	// Специальная обработка корневого маршрута
	if path == "/" || path == "" {
		for _, method := range m.publicRoutes["/"] {
			if method == "*" || ctx.Method() == method {
				return true
			}
		}
	}

	// Удаление последнего слеша
	path = strings.TrimRight(path, "/")

	// Проверка других маршрутов
	for route, methods := range m.publicRoutes {
		if utils.IsRouteMatch(path, route) {
			for _, method := range methods {
				if method == "*" || ctx.Method() == method {
					return true
				}
			}
		}
	}
	return false
}

// extractTokenClaims извлекает и проверяет токен из заголовка Authorization.
func (m *AuthMiddleware) extractTokenClaims(ctx *fiber.Ctx) (*utils.JwtCustomClaim, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return nil, utils.ErrorEmptyAuth
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, utils.ErrorInvalidFormatToken
	}

	token := parts[1]
	// Проверка токена через сервис JWT
	claims, err := m.jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
