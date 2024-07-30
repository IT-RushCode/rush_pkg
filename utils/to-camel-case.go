package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// toCamelCase переводит snake_case в camelCase
func ToCamelCase(s string) string {
	s = strings.ToLower(s)
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	caser := cases.Title(language.Und)
	for i := 1; i < len(parts); i++ {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}
