package utils

import (
	"errors"
	"strings"
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
)

// Отлов ошибки дубликата и вывод соответствующей ошибки.
func HandleDuplicateKeyError(err error) error {
	if err == nil {
		return nil
	}

	// PostgreSQL
	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return ErrExists
	}

	// MySQL
	if strings.Contains(err.Error(), "Duplicate entry") {
		return ErrExists
	}

	// MSSQL
	if strings.Contains(err.Error(), "Violation of UNIQUE KEY constraint") {
		return ErrExists
	}

	return ErrCreate
}
