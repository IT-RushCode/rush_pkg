package utils

import (
	"fmt"
	"time"
)

const (
	RusDateTimeFormat = "02.01.2006 15:04:05"
	RusTimeDateFormat = "15:04:05 02.01.2006"
	RusDateFormat     = "02.01.2006"
	RusTimeFormat     = "15:04:05"
)

func ParseDateTime(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil
	}

	// Часто используемые форматы дат и времени
	formats := []string{
		time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
		time.DateTime,          // "2006-01-02 15:04:05"
		time.DateOnly,          // "2006-01-02"
		RusDateTimeFormat,      // "02.01.2006 15:04:05"
		RusTimeDateFormat,      // "15:04:05 02.01.2006"
		RusDateFormat,          // "02.01.2006"
		RusTimeFormat,          // "15:04:05"
		"2006-01-02T15:04:05",  // ISO-8601 без зоны
		"02-01-2006 15:04:05",  // Европейский формат с временем
		"02-01-2006",           // Европейский формат без времени
		"02/01/2006 15:04:05",  // Формат с временем с точкой
		"02/01/2006",           // Формат без времени с точкой
		"2006/01/02 15:04:05",  // Год/месяц/день с временем через слеш
		"2006/01/02",           // Год/месяц/день без времени
		"02 Jan 2006 15:04:05", // День месяц год, текстовый месяц
		"02 Jan 2006",          // День месяц год, без времени
		time.RFC1123,           // RFC1123 с временной зоной
	}

	var t time.Time
	var err error
	for _, format := range formats {
		if t, err = time.Parse(format, dateTimeStr); err == nil {
			return &t, nil
		}
	}
	return nil, err
}

// Вспомогательная функция для безопасного форматирования даты и времени
func CheckDateTimeToNil(date *time.Time, format string) string {
	if date == nil || date.IsZero() {
		return "" // Возвращаем пустую строку, если дата отсутствует
	}

	if format == "" {
		format = time.RFC3339
	}

	return date.Format(format) // Форматируем только дату
}

// SetLocalTimezone устанавливает локальную временную зону на основе переданного имени зоны
func SetLocalTimezone(timezone string) error {
	if timezone == "" {
		return fmt.Errorf("временная зона не указана")
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке временной зоны: %v", err)
	}

	// Устанавливаем локальную временную зону для всего приложения
	time.Local = loc
	return nil
}
