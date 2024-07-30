package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

func SendResponse(ctx *fiber.Ctx, success bool, message string, body interface{}, statusCode int) error {
	response := Response{
		Status:  success,
		Message: message,
		Body:    body,
	}

	ctx.Status(statusCode)

	if err := ctx.JSON(response); err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return nil
}

// SuccessResponse отправляет успешный ответ с указанным сообщением и телом данных.
func SuccessResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, true, message, body, fiber.StatusOK)
}

// ErrorResponse отправляет ответ об ошибке сервера с указанным сообщением и телом данных.
func ErrorResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusInternalServerError)
}

// ErrorNotFoundResponse отправляет ответ об ошибке "Не найдено" с указанным сообщением и телом данных.
func ErrorNotFoundResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusNotFound)
}

// ErrorBadRequestResponse отправляет ответ об ошибке "Неверный запрос" с указанным сообщением и телом данных.
func ErrorBadRequestResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusBadRequest)
}

// ErrorUnauthorizedResponse отправляет ответ об ошибке "Неавторизован" с указанным сообщением и телом данных.
func ErrorUnauthorizedResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusUnauthorized)
}

// ErrorForbiddenResponse отправляет ответ об ошибке "Запрещено" с указанным сообщением и телом данных.
func ErrorForbiddenResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusForbidden)
}

// ErrorConflictResponse отправляет ответ об ошибке "Конфликт" с указанным сообщением и телом данных.
func ErrorConflictResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusConflict)
}

// ErrorUnsupportedMediaTypeResponse отправляет ответ об ошибке "Неподдерживаемый тип медиа" с указанным сообщением и телом данных.
func ErrorUnsupportedMediaTypeResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusUnsupportedMediaType)
}

// ErrorInternalServerErrorResponse отправляет ответ об ошибке "Внутренняя ошибка сервера" с указанным сообщением и телом данных.
func ErrorInternalServerErrorResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusInternalServerError)
}

var (
	Success = "успешно"
)
