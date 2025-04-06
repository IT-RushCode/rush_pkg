package middlewares

import (
	"net/http"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorMiddleware struct {
	env string
	log *zap.Logger
}

// NewErrorMiddleware создает новый экземпляр ErrorMiddleware.
func NewErrorMiddleware(cfg *config.AppConfig, log *zap.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{
		env: cfg.ENV,
		log: log,
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
	if m.env == "dev" {
		// Логируем детализированную ошибку для отладки
		m.log.Debug("Ошибка в запросе", utils.WithRID(ctx), zap.Error(cleanedError))

		// Возвращаем детализированное сообщение клиенту с правильным HTTP статусом
		return utils.SendResponse(ctx, false, cleanedError.Error(), nil, statusCode)
	}

	// В Production режиме скрываем детали ошибки и используем клиентские сообщения
	clientMessage := utils.GetClientErrorMessage(cleanedError)

	// Если в режиме production, логируем ошибку
	if m.env == "prod" {
		// Логируем ошибку с подробностями в режиме production
		m.log.Error("Внутренняя ошибка", utils.WithRID(ctx), zap.Error(cleanedError))
	}

	if clientMessage == utils.ErrClientInternal {
		statusCode = http.StatusInternalServerError
	}

	// Возвращаем клиенту обобщённое или специфическое сообщение об ошибке
	return utils.SendResponse(ctx, false, clientMessage, nil, statusCode)
}
