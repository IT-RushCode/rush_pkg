package middlewares

import (
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// PermissionChecker определяет интерфейс для проверки прав пользователя.
type PermissionChecker interface {
	// HasPermission проверяет, есть ли у пользователя определенная привилегия.
	HasPermission(userID uint, permission string) bool
}

// PermissionMiddleware представляет middleware для проверки привилегий пользователя.
type PermissionMiddleware struct {
	// checker реализует интерфейс PermissionChecker для проверки прав доступа.
	checker PermissionChecker
}

// NewPermissionMiddleware создает новый экземпляр PermissionMiddleware с заданным проверяющим.
func NewPermissionMiddleware(checker PermissionChecker) *PermissionMiddleware {
	return &PermissionMiddleware{
		checker: checker,
	}
}

// CheckPermission возвращает обработчик middleware, который проверяет наличие у пользователя требуемой привилегии.
func (p *PermissionMiddleware) CheckPermission(requiredPermission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Проверяем, является ли маршрут публичным.
		isPublic, _ := ctx.Locals("IsPublic").(bool)
		if isPublic {
			// Если маршрут публичный, пропускаем проверку привилегий.
			return ctx.Next()
		}

		// Получаем UserID из локальных данных контекста.
		userId, ok := ctx.Locals("UserID").(uint)
		if !ok {
			// Если UserID отсутствует или неверного типа, возвращаем ошибку доступа.
			return utils.ErrorForbiddenResponse(ctx, "неверный формат UserID", nil)
		}

		// Проверяем, имеет ли пользователь необходимую привилегию.
		if !p.checker.HasPermission(userId, requiredPermission) && requiredPermission != "me" {
			// Если привилегии недостаточно, возвращаем ошибку доступа.
			return utils.ErrorForbiddenResponse(ctx, utils.ErrForbidden.Error(), nil)
		}

		// Если все проверки пройдены, передаем управление следующему обработчику.
		return ctx.Next()
	}
}
