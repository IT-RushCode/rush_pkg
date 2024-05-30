package utils

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrId              = errors.New("неверный id")
	ErrCreate          = errors.New("ошибка создания записи")
	ErrUpdate          = errors.New("ошибка обновления записи")
	ErrDelete          = errors.New("ошибка удаления записи")
	ErrExists          = errors.New("запись с такими же параметрами уже существует")
	ErrRecordNotFound  = errors.New("запись не найдена")
	ErrRecordsNotFound = errors.New("записи не найдены")
	ErrFileNotFound    = errors.New("файл не найден")
	ErrFilesNotFound   = errors.New("файлы не найдены")
	ErrPermission      = errors.New("нет прав на редактирование")
	ErrDuplicate       = errors.New("дубликат записи")
	ErrInvalidInput    = errors.New("ошибка входящих данных")
	ErrInternal        = errors.New("внутренняя ошибка сервера")
)

func HandleDuplicateKeyError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
		strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "Violation of UNIQUE KEY constraint") {
		return ErrExists
	}

	return ErrCreate
}

// Функция для проверки возвращаемой ошибки из репозитория
func CheckErr(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrRecordNotFound):
		return ErrorNotFoundResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrExists):
		return ErrorConflictResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrInvalidInput):
		return ErrorBadRequestResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrPermission):
		return ErrorForbiddenResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrInternal):
		return ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	default:
		return ErrorResponse(ctx, err.Error(), nil)
	}
}
