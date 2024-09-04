package utils

// Ternary реализует аналог тернарного оператора
//
// Пример использования: max := Ternary(a > b, a, b)
func Ternary[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
