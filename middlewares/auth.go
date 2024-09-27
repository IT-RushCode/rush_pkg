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
	jwtTTL       time.Duration
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
	jwtService := utils.NewJWTService(cfg.JWT.JWT_SECRET, jwtTTL, jwtRTTL)

	return &AuthMiddleware{
		jwtTTL:       jwtTTL,
		jwt:          jwtService,
		publicRoutes: routes,
	}
}

// VerifyToken выполняет основную проверку токена.
func (m *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	// Проверка на публичный список маршрутов
	if m.isPublicRoute(ctx) {
		ctx.Locals("IsPublic", true)
		return ctx.Next()
	}

	// Проверка и извлечение данных из токена
	claims, err := m.extractTokenClaims(ctx)
	if claims == nil {
		return err
	}

	// Сохранение данных пользователя в контексте
	ctx.Locals("UserID", claims.UserID)
	ctx.Locals("IsMob", claims.IsMob)

	// Если привилегии проверены, выполняем следующего обработчика
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
