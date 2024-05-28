package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func ValidateStruct(data interface{}) error {
	// Регистрация кастомных валидаторов
	if err := registerCustomValidators(validate); err != nil {
		return fmt.Errorf("ошибка при регистрации кастомных валидаторов: %v", err)
	}

	err := validate.Struct(data)
	if err != nil {
		var errMsg strings.Builder

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("валидация не удалась: %s", err)
		}

		errMsg.WriteString("Ошибка валидации: ")

		errors := err.(validator.ValidationErrors)
		for i, err := range errors {
			errMsg.WriteString(fmt.Sprintf("'%s' не прошло валидацию по правилу '%s'", err.StructField(), err.Tag()))
			if err.Param() != "" {
				errMsg.WriteString(fmt.Sprintf(" = %+v", err.Param()))
			}

			if i < len(errors)-1 {
				errMsg.WriteString("; ")
			} else {
				errMsg.WriteString(".")
			}
		}

		return fmt.Errorf("%s", errMsg.String())
	}

	return nil
}

// validatePhone проверяет формат телефонного номера
func validatePhone(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Регулярное выражение для проверки формата телефонного номера
	re := regexp.MustCompile(`^\+\d{1,3}\d{3}\d{3}\d{2}\d{2}$`)
	return re.MatchString(phoneNumber)
}

// registerCustomValidators регистрирует кастомные валидаторы
func registerCustomValidators(validate *validator.Validate) error {
	// Регистрируем кастомный валидатор для формата телефонного номера
	if err := validate.RegisterValidation("phone", validatePhone); err != nil {
		return err
	}
	return nil
}
