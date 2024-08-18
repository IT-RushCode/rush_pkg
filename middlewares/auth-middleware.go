package middlewares

import (
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthData struct {
	UserID uint   `json:"userId"`
	IP     string `json:"ip"`
}

// AuthMiddleware представляет собой middleware для аутентификации пользователя.
type AuthMiddleware struct {
	jwtTTL          time.Duration
	jwt             utils.JWTService
	whiteListRoutes map[string][]string
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(cfg *config.Config, routes map[string][]string) *AuthMiddleware {
	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Second
	jwtRTTL := time.Duration(cfg.JWT.REFRESH_TTL) * time.Second
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL)

	return &AuthMiddleware{
		jwtTTL:          jwtTTL,
		jwt:             jwtService,
		whiteListRoutes: routes,
	}
}

// VerifyToken проверяет токен аутентификации пользователя.
func (m *AuthMiddleware) VerifyToken(ctx *fiber.Ctx) error {
	// Проверка маршрута и метода в бело		м списке
	for route, methods := range m.whiteListRoutes {
		if isRouteMatch(ctx.Path(), route) {
			for _, method := range methods {
				if method == "*" || ctx.Method() == method {
					// Если маршрут и метод в белом списке, или метод указан как "*", пропускаем проверку токена
					return ctx.Next()
				}
			}
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

// isRouteMatch проверяет, соответствует ли путь маршруту, включая поддержку параметров.
func isRouteMatch(path, route string) bool {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false
	}

	for i := range routeParts {
		if routeParts[i] != pathParts[i] && !strings.HasPrefix(routeParts[i], ":") {
			return false
		}
	}
	return true
}
