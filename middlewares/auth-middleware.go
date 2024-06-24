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
		"/api/v1/auth/phone-login",
		"/api/v1/auth/refresh-token",
		// "/swagger",
		"/fb-metrics",
		"/metrics",
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

// // AuthorizationMiddleware проверяет права пользователя
// func AuthorizationMiddleware(requiredPermission string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Предполагается, что вы уже получили userID из токена или сессии
// 		userID := c.Locals("userID")
// 		if userID == nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
// 		}

// 		var user *rpModels.User
// 		if err := db.Preload("Permissions").Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не найден"})
// 		}

// 		// Проверка прав пользователя
// 		hasPermission := false
// 		for _, perm := range user.Permissions {
// 			if perm.Name == requiredPermission {
// 				hasPermission = true
// 				break
// 			}
// 		}

// 		// Проверка прав через роли
// 		if !hasPermission {
// 			for _, role := range user.Roles {
// 				for _, perm := range role.Permissions {
// 					if perm.Name == requiredPermission {
// 						hasPermission = true
// 						break
// 					}
// 				}
// 			}
// 		}

// 		if !hasPermission {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Доступ запрещен"})
// 		}

// 		return c.Next()
// 	}
// }
