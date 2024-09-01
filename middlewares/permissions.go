package middlewares

import (
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// PermissionChecker определяет метод для проверки прав пользователя.
type PermissionChecker interface {
	HasPermission(userID uint, permission string) bool
}

// CheckPermissions проверяет, есть ли у пользователя привилегия для текущего маршрута.
func (m *AuthMiddleware) CheckPermissions(ctx *fiber.Ctx, claims *utils.JwtCustomClaim) error {
	routeName := ctx.Route().Name
	if routeName == "" {
		return utils.ErrorForbiddenResponse(ctx, "маршрут не имеет привилегии", nil)
	}

	// Проверка привилегий через интерфейс
	if !m.permissionChecker.HasPermission(claims.UserID, routeName) {
		return utils.ErrorForbiddenResponse(ctx, "доступ запрещен", nil)
	}

	return nil
}
