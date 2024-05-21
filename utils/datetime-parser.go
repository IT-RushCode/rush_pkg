package utils

import "time"

func ParseDateTime(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil
	}

	// Определяем формат времени ISO 8601 с миллисекундами и Z
	const layout = "2006-01-02T15:04:05.999Z"
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return nil, err
	}

	// Применяем временную зону "+0500"
	loc, _ := time.LoadLocation("")  // Создаем локацию для UTC+5
	t = t.In(loc).Add(5 * time.Hour) // Прибавляем 5 часов к UTC

	return &t, nil
}
