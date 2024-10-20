package middlewares

import (
	"log"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ErrorMiddleware struct {
	dbg bool
}

// NewErrorMiddleware создает новый экземпляр ErrorMiddleware.
func NewErrorMiddleware(cfg *config.AppConfig) *ErrorMiddleware {
	return &ErrorMiddleware{
		dbg: cfg.DEBUG,
	}
}

// ErrorHandlingMiddleware обрабатывает ошибки и возвращает правильные сообщения в зависимости от режима (Debug или Production)
func (m *ErrorMiddleware) ErrorHandlingMiddleware(ctx *fiber.Ctx) error {
	// Обрабатываем запрос, если есть ошибка - перехватываем её
	err := ctx.Next()

	// Проверка на наличие кастомной ошибки, которая может обходить middleware
	if customErr, ok := ctx.Locals("CustomError").(utils.CustomErrorRes); ok {
		return utils.SendResponse(ctx, false, customErr.Message, customErr.Body, customErr.StatusCode)
	}

	// Проверка на наличие кастомного кода, который может обходить middleware
	if _, ok := ctx.Locals("CustomCode").(bool); ok {
		// Возвращаем оригинальную ошибку, обработка уже произошла
		return err
	}

	// Если нет ошибки, продолжаем
	if err == nil {
		return nil
	}

	// Пропускаем ошибку через MapErrorToStatus для получения правильного HTTP статуса и клиентского сообщения
	statusCode, cleanedError := utils.MapErrorToStatus(err)

	// В режиме Debug возвращаем детализированную ошибку
	if m.dbg {
		// Логируем детализированную ошибку для отладки
		log.Printf("%s(DEBUG)%s Error > %s %v", utils.Green, utils.Red, cleanedError, utils.Reset)

		// Возвращаем детализированное сообщение клиенту с правильным HTTP статусом
		return utils.SendResponse(ctx, false, cleanedError.Error(), nil, statusCode)
	}

	// В Production режиме скрываем детали ошибки и используем клиентские сообщения
	clientMessage := utils.GetClientErrorMessage(cleanedError)

	// Логируем внутреннюю ошибку для мониторинга в Production
	log.Printf("%s(PROD)%s Внутренняя ошибка > %s %v", utils.Green, utils.Red, cleanedError, utils.Reset)

	// Возвращаем клиенту обобщённое или специфическое сообщение об ошибке
	return utils.SendResponse(ctx, false, clientMessage, nil, statusCode)
}
