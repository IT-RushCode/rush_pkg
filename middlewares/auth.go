package middlewares

import (
	"fmt"
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

// PermissionChecker определяет метод для проверки прав пользователя.
type PermissionChecker interface {
	HasPermission(userID uint, permission string) bool
}

// AuthMiddleware представляет собой middleware для аутентификации пользователя.
type AuthMiddleware struct {
	jwtTTL            time.Duration
	jwt               utils.JWTService
	publicRoutes      map[string][]string
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
		jwtTTL:            jwtTTL,
		jwt:               jwtService,
		publicRoutes:      routes,
		permissionChecker: checker,
	}
}

// VerifyToken выполняет основную проверку токена и привилегий.
func (m *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	// Удаление последнего слеша
	normalizePath(ctx)

	// Проверка на публичный список маршрутов
	if m.isWhiteListedRoute(ctx) {
		return ctx.Next()
	}

	// Проверка и извлечение данных из токена
	claims, err := m.extractTokenClaims(ctx)
	if claims == nil {
		return err
	}

	// Сохранение данных пользователя в контексте
	ctx.Locals("UserID", claims.UserID)

	// Вызов `ctx.Next()` для перехода к следующему обработчику маршрута
	if err := ctx.Next(); err != nil {
		return err
	}

	// Проверка привилегий пользователя на основе имени маршрута
	if err := m.checkPermissions(ctx, claims); err != nil {
		return err
	}

	return nil
}

// isWhiteListedRoute проверяет, является ли маршрут и метод запросом в белом списке.
func (m *AuthMiddleware) isWhiteListedRoute(ctx *fiber.Ctx) bool {
	for route, methods := range m.publicRoutes {
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

// CheckPermissions проверяет, есть ли у пользователя привилегия для те	кущего маршрута.
func (m *AuthMiddleware) checkPermissions(ctx *fiber.Ctx, claims *utils.JwtCustomClaim) error {
	routeName := ctx.Route().Name
	if routeName == "" {
		return utils.ErrorForbiddenResponse(ctx, fmt.Sprintf("маршрут %s не имеет привилегии", ctx.Route().Path), nil)
	}

	if routeName == "me" {
		return nil
	}

	// Проверка привилегий через интерфейс PermissionChecker
	if !m.permissionChecker.HasPermission(claims.UserID, routeName) {
		return utils.ErrorForbiddenResponse(ctx, "доступ запрещен", nil)
	}

	return nil
}

// normalizePath нормализует путь, удаляя последний слеш, если он есть.
func normalizePath(ctx *fiber.Ctx) {
	ctx.Path(strings.TrimRight(ctx.Path(), "/"))
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
