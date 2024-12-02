package utils

// OrDefault возвращает значение указателя или значение по умолчанию, если указатель nil.
func OrDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
