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
	jwtTTL            time.Duration
	jwt               utils.JWTService
	whiteListRoutes   map[string][]string
	permissionChecker PermissionChecker
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(
	cfg *config.Config,
	routes map[string][]string,
	checker PermissionChecker,
) *AuthMiddleware {
	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Second
	jwtRTTL := time.Duration(cfg.JWT.REFRESH_TTL) * time.Second
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL)

	return &AuthMiddleware{
		jwtTTL:          jwtTTL,
		jwt:             jwtService,
		whiteListRoutes: routes,
	}
}

// VerifyToken выполняет основную проверку токена и привилегий.
func (m *AuthMiddleware) VerifyToken(ctx *fiber.Ctx) error {
	// Удаление последнего слеша
	normalizePath(ctx)

	// Проверка на белый список маршрутов
	if m.isWhiteListedRoute(ctx) {
		return ctx.Next()
	}

	// Проверка и извлечение данных из токена
	claims, err := m.extractTokenClaims(ctx)
	if err != nil {
		return err
	}

	// Сохранение данных пользователя в контексте
	ctx.Locals("UserID", claims.UserID)

	// Проверка привилегий пользователя на основе имени маршрута
	if err := m.CheckPermissions(ctx, claims); err != nil {
		return err
	}

	return ctx.Next()
}

// normalizePath нормализует путь, удаляя последний слеш, если он есть.
func normalizePath(ctx *fiber.Ctx) {
	ctx.Path(strings.TrimRight(ctx.Path(), "/"))
}

// isWhiteListedRoute проверяет, является ли маршрут и метод запросом в белом списке.
func (m *AuthMiddleware) isWhiteListedRoute(ctx *fiber.Ctx) bool {
	for route, methods := range m.whiteListRoutes {
		if isRouteMatch(ctx.Path(), route) {
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
		return nil, utils.ErrorUnauthorizedResponse(ctx, "отсутствует токен авторизации", nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, utils.ErrorUnauthorizedResponse(ctx, "неверный формат токена", nil)
	}

	token := parts[1]

	// Проверка токена через сервис JWT
	claims, err := m.jwt.ValidateToken(token)
	if err != nil {
		return nil, utils.ErrorUnauthorizedResponse(ctx, "неверный токен авторизации", nil)
	}

	return claims, nil
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
