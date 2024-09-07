package middlewares

import (
	"fmt"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// PermissionChecker определяет метод для проверки прав пользователя.
type PermissionChecker interface {
	HasPermission(userID uint, permission string) bool
}

type PermissionsMiddleware struct {
	checker PermissionChecker
}

func NewPermissionsMiddleware(checker PermissionChecker) *PermissionsMiddleware {
	return &PermissionsMiddleware{checker: checker}
}

func (p *PermissionsMiddleware) CheckPermissions(ctx *fiber.Ctx) error {
	// Получаем имя маршрута
	route := ctx.Route()

	if route.Name == "" {
		return utils.ErrorForbiddenResponse(ctx, fmt.Sprintf("маршрут %s не имеет привилегии", route.Path), nil)
	}

	// Получаем UserID из локальных данных как int
	userId, ok := ctx.Locals("UserID").(uint)
	if !ok {
		return utils.ErrorForbiddenResponse(ctx, "неверный формат UserID", nil)
	}

	// Проверка прав доступа пользователя
	if !p.checker.HasPermission(userId, route.Name) {
		// Если привилегии не соответствуют, возвращаем ошибку и не продолжаем выполнение
		return utils.ErrorForbiddenResponse(ctx, utils.ErrForbidden.Error(), nil)
	}

	return ctx.Next()
}
