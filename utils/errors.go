package utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	// Repositoriy errors R1000-R1999
	ErrCreate             = errors.New("RP1000: Ошибка создания записи")
	ErrUpdate             = errors.New("RP1001: Ошибка обновления записи")
	ErrDelete             = errors.New("RP1002: Ошибка удаления записи")
	ErrExists             = errors.New("RP1003: Запись с такими же параметрами уже существует")
	ErrRecordNotFound     = errors.New("RP1004: Запись не найдена")
	ErrRecordsNotFound    = errors.New("RP1005: Записи не найдены")
	ErrInvalidTransaction = errors.New("RP1006: Некорректная транзакция")

	// Auth errors A2000-A2999
	ErrCopierData                = errors.New("AU2000: Ошибка copier")
	ErrRefreshToken              = errors.New("AU2001: Неверный токен обновления")
	ErrNotRefreshToken           = errors.New("AU2002: Полученный токен не является refresh токеном")
	ErrUserIdNotFound            = errors.New("AU2003: User id не найден в context")
	ErrorEmptyAuth               = errors.New("AU2004: Отсутствует токен в запросе")
	ErrForbidden                 = errors.New("AU2005: Нет прав")
	ErrUnauthenticated           = errors.New("AU2006: Не авторизован")
	ErrorGenAccessToken          = errors.New("AU2007: Не удалось сгенерировать токен доступа")
	ErrorGenRefreshToken         = errors.New("AU2008: Не удалось создать токен обновления")
	ErrorSigningMethod           = errors.New("AU2009: Неверный метод подписи токена")
	ErrorInvalidToken            = errors.New("AU2010: Неверный токен")
	ErrorInvalidFormatToken      = errors.New("AU2011: Неверный токен")
	ErrorTokenExpired            = errors.New("AU2012: Время действия токена истекло")
	ErrorTokenNotYetValid        = errors.New("AU2013: Токен недействителен")
	ErrorRefreshTokenExpired     = errors.New("AU2014: Время действия токена обновления истекло")
	ErrorRefreshTokenNotYetValid = errors.New("AU2015: Токен обновления недействителен")

	// BadRequestErrors 4000-4999
	ErrInvalidData  = errors.New("BR4000: Некорректные данные")
	ErrInvalidBody  = errors.New("BR4001: Неверный формат тела запроса")
	ErrInvalidField = errors.New("BR4002: Некорректное поле в запросе")
	ErrGetUUID      = errors.New("BR4003: Не указан uuid")
	ErrInvalidUUID  = errors.New("BR4004: Неверный формат UUID")

	// File handler errors 5000-5999
	ErrUploadFile     = errors.New("FL5000: Ошибка загрузки файла")
	ErrUpdateFile     = errors.New("FL5001: Ошибка обновления файла")
	ErrDeleteFile     = errors.New("FL5002: Ошибка удаления файла")
	ErrDeleteOldFile  = errors.New("FL5003: Ошибка удаления старого файла")
	ErrSaveFile       = errors.New("FL5004: Ошибка сохранения файла")
	ErrFileNotFound   = errors.New("FL5005: Файл не найден")
	ErrFilesNotFound  = errors.New("FL5006: Файлы не найдены")
	ErrCreateDir      = errors.New("FL5007: Ошибка создания директории")
	ErrRecordImageble = errors.New("FL5008: Запись с указанным imagebleID не существует")
	ErrSaveMetaData   = errors.New("FL5009: Ошибка сохранения метаданных файла")
	ErrUpdateMetaData = errors.New("FL5010: Ошибка обновления метаданных файла")
	ErrDeleteMetaData = errors.New("FL5011: Ошибка удаления метаданных файла")

	// Context errors 6000-6999
	ErrDeadlineExceeded = errors.New("CTX6000: Превышен тайм-аут операции")
	ErrCancelContext    = errors.New("CTX6001: Операция была отменена")

	// Internal errors ST500
	ErrInternal       = errors.New("ST500: Внутренняя ошибка сервера")
	ErrClientInternal = "Что-то пошло не так. Пожалуйста, повторите попытку позже."
)

// Функция для проверки возвращаемой ошибки из репозитория
func CheckErr(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFoundResponse(ctx, err.Error(), nil)
	case strings.Contains(err.Error(), ErrExists.Error()):
		return ErrorConflictResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrForbidden):
		return ErrorForbiddenResponse(ctx, err.Error(), nil)
	case errors.Is(err, ErrInternal):
		return ErrorInternalServerErrorResponse(ctx, err.Error(), nil)
	default:
		return ErrorResponse(ctx, err.Error(), nil)
	}
}

func MapErrorToStatus(err error) (int, error) {
	// Проверка на доменные ошибки
	switch err {
	case ErrCreate:
		return http.StatusInternalServerError, ErrCreate
	case ErrUpdate:
		return http.StatusInternalServerError, ErrUpdate
	case ErrDelete:
		return http.StatusInternalServerError, ErrDelete
	case ErrExists:
		return http.StatusConflict, ErrExists
	case ErrRecordNotFound:
		return http.StatusNotFound, ErrRecordNotFound
	case ErrRecordsNotFound:
		return http.StatusNotFound, ErrRecordsNotFound
	case ErrInternal:
		return http.StatusInternalServerError, ErrInternal
	case ErrForbidden:
		return http.StatusForbidden, ErrForbidden
	case ErrUnauthenticated:
		return http.StatusUnauthorized, ErrUnauthenticated
	}

	// Обработка ошибок GORM (например, ошибки базы данных)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, ErrRecordNotFound
	}

	if errors.Is(err, gorm.ErrInvalidData) {
		return http.StatusBadRequest, ErrInvalidData
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return http.StatusConflict, err
	}

	if errors.Is(err, gorm.ErrInvalidTransaction) {
		return http.StatusConflict, ErrInvalidTransaction
	}

	if errors.Is(err, gorm.ErrInvalidField) {
		return http.StatusBadRequest, ErrInvalidField
	}

	// Обработка ошибок контекста (тайм-ауты или отмены)
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, ErrDeadlineExceeded
	}

	if errors.Is(err, context.Canceled) {
		return http.StatusRequestTimeout, ErrCancelContext
	}

	if strings.Contains(err.Error(), "Unprocessable Entity") {
		return http.StatusUnprocessableEntity, ErrInvalidBody
	}

	// Проверка на ошибки валидации по префиксу "VALIDATE:"
	if strings.HasPrefix(err.Error(), "VALIDATE:") {
		log.Println("Validation error:", err)

		// Удаляем префикс "VALIDATE:" перед выводом ошибки
		return http.StatusBadRequest, errors.New(strings.TrimPrefix(err.Error(), "VALIDATE: "))
	}

	if errors.Is(err, ErrorEmptyAuth) {
		return http.StatusUnauthorized, ErrorEmptyAuth
	}

	// Обработка ошибок декодирования токена
	if errors.Is(err, ErrorInvalidFormatToken) || strings.Contains(err.Error(), "token contains an invalid number of segments") ||
		strings.Contains(err.Error(), "could not base64 decode") ||
		strings.Contains(err.Error(), "could not JSON decode") {
		return http.StatusBadRequest, ErrorInvalidFormatToken
	}

	// Обработка ошибки подписи токена
	if strings.Contains(err.Error(), "signature is invalid") {
		return http.StatusUnauthorized, ErrorSigningMethod
	}

	// Проверка на ошибки истечения срока действия токена
	if strings.Contains(err.Error(), "token is expired") {
		return http.StatusUnauthorized, ErrorTokenExpired
	}

	// Если ошибка неизвестного типа или internal, возвращаем внутреннюю ошибку
	if strings.Contains(err.Error(), "internal error") {
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusInternalServerError, err
}

// -------------------- MIDDLEWARE ERRORS -------------------->

var clientErrorMessages = map[error]string{
	// Repository errors
	ErrCreate:             "Ошибка создания записи. Пожалуйста, повторите попытку позже.",
	ErrUpdate:             "Ошибка обновления записи. Пожалуйста, повторите попытку позже.",
	ErrDelete:             "Ошибка удаления записи. Пожалуйста, повторите попытку позже.",
	ErrExists:             "Запись с такими параметрами уже существует.",
	ErrRecordNotFound:     "Запись не найдена.",
	ErrRecordsNotFound:    "Записи не найдены.",
	ErrInvalidTransaction: "Неверная транзакция. Пожалуйста, попробуйте снова.",

	// Auth errors
	ErrCopierData:                "Ошибка обработки данных.",
	ErrRefreshToken:              "Неверный токен обновления. Пожалуйста, авторизуйтесь заново.",
	ErrNotRefreshToken:           "Полученный токен не является токеном обновления.",
	ErrUserIdNotFound:            "Не удалось найти идентификатор пользователя.",
	ErrorEmptyAuth:               "Запрос не содержит данных для авторизации.",
	ErrForbidden:                 "У вас нет прав для выполнения этого действия.",
	ErrUnauthenticated:           "Не авторизован.",
	ErrorGenAccessToken:          "Ошибка при создании токена доступа.",
	ErrorGenRefreshToken:         "Ошибка при создании токена обновления.",
	ErrorSigningMethod:           "Неверный метод подписи токена.",
	ErrorInvalidToken:            "Неверный токен.",
	ErrorTokenExpired:            "Срок действия сессии истёк. Войдите снова.",
	ErrorTokenNotYetValid:        "Доступ невозможен. Попробуйте ещё раз.",
	ErrorRefreshTokenExpired:     "Сессия завершена. Войдите снова.",
	ErrorRefreshTokenNotYetValid: "Сессия недействительна. Войдите снова.",

	// BadRequest errors
	ErrInvalidData:  "Некорректные данные. Проверьте правильность ввода.",
	ErrInvalidBody:  "Некорректные данные. Тело запроса не может быть пустым.",
	ErrInvalidField: "Некорректное поле в запросе.",
	ErrGetUUID:      "UUID не указан.",
	ErrInvalidUUID:  "Неверный формат UUID.",

	// File handler errors
	ErrUploadFile:     "Ошибка загрузки файла. Попробуйте снова.",
	ErrUpdateFile:     "Ошибка обновления файла.",
	ErrDeleteFile:     "Ошибка удаления файла.",
	ErrDeleteOldFile:  "Ошибка удаления старого файла.",
	ErrSaveFile:       "Ошибка сохранения файла.",
	ErrFileNotFound:   "Файл не найден.",
	ErrFilesNotFound:  "Файлы не найдены.",
	ErrCreateDir:      "Ошибка создания директории.",
	ErrRecordImageble: "Запись с указанным идентификатором не найдена.",
	ErrSaveMetaData:   "Ошибка сохранения метаданных файла.",
	ErrUpdateMetaData: "Ошибка обновления метаданных файла.",
	ErrDeleteMetaData: "Ошибка удаления метаданных файла.",

	// Context errors
	ErrDeadlineExceeded: "Превышен тайм-аут операции. Попробуйте позже.",
	ErrCancelContext:    "Операция была отменена.",

	// Internal errors
	ErrInternal: ErrClientInternal,
}

// GetClientErrorMessage возвращает клиентский текст ошибки для Production режима
func GetClientErrorMessage(err error) string {
	if message, exists := clientErrorMessages[err]; exists {
		return message
	}
	// Если ошибка не найдена в списке, возвращаем обобщённое сообщение
	return ErrClientInternal
}
