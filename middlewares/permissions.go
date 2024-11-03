package middlewares

import (
	"context"
	"strconv"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// Основной интерфейс PermissionChecker с обязательными методами.
// PermissionChecker определяет интерфейс репозитория для проверки прав пользователя.
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID uint, permission string) bool
	IsUserActive(ctx context.Context, userID uint) bool
}

// Расширенный интерфейс, добавляющий опциональный метод HasAccessToObject.
type OptionalPermissionChecker interface {
	PermissionChecker
	HasAccessToObject(ctx context.Context, userID, objectID uint, objectType string) (bool, error)
}

// --------------------- PERMISSION MIDDLEWARE --------------------->

// PermissionMiddleware представляет middleware для проверки привилегий пользователя.
type PermissionMiddleware struct {
	checker PermissionChecker // Реализация интерфейса PermissionChecker
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
			return ctx.Next()
		}

		// Получаем UserID из локальных данных контекста.
		userID, ok := ctx.Locals("UserID").(uint)
		if !ok {
			return utils.ErrorForbiddenResponse(ctx, "неверный формат UserID", nil)
		}

		// Проверка статуса пользователя
		if !p.checker.IsUserActive(ctx.Context(), userID) {
			return utils.ErrorUnauthorizedResponse(ctx, utils.ErrUnauthenticated.Error(), nil)
		}

		// Проверяем привилегии пользователя
		if !p.checker.HasPermission(ctx.Context(), userID, requiredPermission) && requiredPermission != "me" {
			return utils.ErrorForbiddenResponse(ctx, utils.ErrForbidden.Error(), nil)
		}

		return ctx.Next()
	}
}

// CheckObjectPermission проверяет доступ к объекту, если метод HasAccessToObject реализован.
func (p *PermissionMiddleware) CheckObjectPermission(objectType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Проверяем, является ли маршрут публичным.
		isPublic, _ := ctx.Locals("IsPublic").(bool)
		if isPublic {
			return ctx.Next()
		}

		userID, ok := ctx.Locals("UserID").(uint)
		if !ok {
			return utils.ErrorForbiddenResponse(ctx, "неверный формат UserID", nil)
		}

		objectID, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return utils.ErrorBadRequestResponse(ctx, "неверный ID объекта", nil)
		}

		// Проверка на реализацию расширенного интерфейса.
		if checkerWithAccess, ok := p.checker.(OptionalPermissionChecker); ok {
			// Выполняем проверку доступа к объекту только если метод реализован.
			hasAccess, err := checkerWithAccess.HasAccessToObject(ctx.Context(), userID, uint(objectID), objectType)
			if err != nil || !hasAccess {
				return utils.ErrorForbiddenResponse(ctx, "доступ к объекту запрещен", nil)
			}
		}

		return ctx.Next()
	}
}
