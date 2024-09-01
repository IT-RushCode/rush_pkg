package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	// Repositoriy errors
	ErrId                = errors.New("неверный id")
	ErrCreate            = errors.New("ошибка создания записи")
	ErrUpdate            = errors.New("ошибка обновления записи")
	ErrDelete            = errors.New("ошибка удаления записи")
	ErrExists            = errors.New("запись с такими же параметрами уже существует")
	ErrRecordNotFound    = errors.New("запись не найдена")
	ErrRecordsNotFound   = errors.New("записи не найдены")
	ErrDuplicate         = errors.New("дубликат записи")
	ErrFieldNotSupported = "фильтрация по полю '%s' не поддерживается"

	// JWT errors
	ErrorGenAccessToken   = errors.New("не удалось сгенерировать токен доступа")
	ErrorGenRefreshToken  = errors.New("не удалось создать токен обновления")
	ErrorSigningMethod    = errors.New("неверный метод подписи токена")
	ErrorInvalidToken     = errors.New("неверный токен")
	ErrorTokenExpired     = errors.New("токен истёк")
	ErrorTokenNotYetValid = errors.New("токен больше не валидный")
	ErrRefreshToken       = errors.New("неверный токен обновления")
	ErrNotRefreshToken    = errors.New("полученный токен не является refresh токеном")

	// Global errors
	ErrInternal       = errors.New("внутренняя ошибка сервера")
	ErrForbidden      = errors.New("доступ запрещен")
	ErrRecordImageble = errors.New("запись с указанным imagebleID не существует")

	// File handler errors
	ErrUploadFile    = errors.New("ошибка загрузки файла")
	ErrUpdateFile    = errors.New("ошибка обновления файла")
	ErrDeleteFile    = errors.New("ошибка удаления файла")
	ErrDeleteOldFile = errors.New("ошибка удаления старого файла")
	ErrSaveFile      = errors.New("ошибка сохранения файла")
	ErrFileNotFound  = errors.New("файл не найден")
	ErrFilesNotFound = errors.New("файлы не найдены")
	ErrCreateDir     = errors.New("ошибка создания директории")

	// File Metadata
	ErrSaveMetaData   = errors.New("ошибка сохранения метаданных файла")
	ErrUpdateMetaData = errors.New("ошибка обновления метаданных файла")
	ErrDeleteMetaData = errors.New("ошибка удаления метаданных файла")

	// UUID errors
	ErrGetUUID     = errors.New("не указан uuid")
	ErrInvalidUUID = errors.New("неверный формат UUID")

	// Controller errors
	ErrorIncorrectID     = errors.New("некорректный :id в параметре пути").Error()
	ErrorIncorrectUserID = errors.New("некорректный :userId в параметре пути").Error()
	ErrInvalidInput      = errors.New("ошибка входящих данных")
)

type DuplicateKeyError struct {
	Field string
	Msg   string
}

func (e *DuplicateKeyError) Error() string {
	return e.Msg
}

func HandleDuplicateKeyError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
		strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "Violation of UNIQUE KEY constraint") {

		// Регулярное выражение для поиска имени уникального ограничения и поля
		re := regexp.MustCompile(`(?i)(duplicate key value violates unique constraint|Duplicate entry|Violation of UNIQUE KEY constraint)\s*"([^"]+)"`)
		match := re.FindStringSubmatch(err.Error())

		if len(match) > 2 {
			uniqueConstraint := match[2]
			return &DuplicateKeyError{
				Field: uniqueConstraint,
				Msg:   ErrExists.Error() + ": " + uniqueConstraint,
			}
		}

		return &DuplicateKeyError{
			Field: "",
			Msg:   ErrExists.Error(),
		}
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
	case strings.Contains(err.Error(), ErrExists.Error()):
		return ErrorConflictResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrInvalidInput):
		return ErrorBadRequestResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrForbidden):
		return ErrorForbiddenResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrInternal):
		return ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	default:
		return ErrorResponse(ctx, err.Error(), nil)
	}
}
