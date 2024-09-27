package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate  = validator.New()
	once      sync.Once // Используем для того, чтобы регистрация произошла только один раз
)

func ValidateStruct(data interface{}) error {
	// Используем sync.Once, чтобы гарантировать однократную регистрацию кастомных валидаторов
	once.Do(func() {
		err := registerCustomValidators(validate)
		if err != nil {
			panic(fmt.Sprintf("ошибка при регистрации кастомных валидаторов: %v", err))
		}
	})

	err := validate.Struct(data)
	if err != nil {
		var errMsg strings.Builder

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("валидация не удалась: %s", err)
		}

		errors := err.(validator.ValidationErrors)
		for i, err := range errors {
			// Получаем имя поля из тега json
			fieldName := getJSONTag(data, err.StructField())

			// Формирование детализированного сообщения об ошибке
			switch err.Tag() {
			case "required":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' обязательно для заполнения", fieldName))
			case "email":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' должно быть валидным email адресом", fieldName))
			case "min":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' должно содержать минимум %s символов", fieldName, err.Param()))
			case "max":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' должно содержать не более %s символов", fieldName, err.Param()))
			case "phone":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' должно быть валидным номером телефона (+7XXXXXXXXXX)", fieldName))
			case "required_if_false":
				errMsg.WriteString(fmt.Sprintf("Поле '%s' обязательно, если поле '%s' имеет значение false", fieldName, err.Param()))
			default:
				errMsg.WriteString(fmt.Sprintf("Поле '%s' не прошло валидацию: %s", fieldName, err.Tag()))
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

// Функция для извлечения имени поля из тега json
func getJSONTag(data interface{}, fieldName string) string {
	// Получаем тип переданной структуры
	val := reflect.TypeOf(data)
	// Проверяем, является ли она указателем, и получаем значение, если это так
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// Получаем поле по имени
	field, ok := val.FieldByName(fieldName)
	if !ok {
		// Если поле не найдено, возвращаем исходное имя структуры
		return fieldName
	}
	// Извлекаем тег json
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	// В случае наличия нескольких тегов (например, `json:"field,omitempty"`), берем первое слово
	return strings.Split(jsonTag, ",")[0]
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

	// Регистрируем кастомный валидатор для логики зависимости полей
	if err := validate.RegisterValidation("required_if_false", validateRequiredIfFalse); err != nil {
		return err
	}

	return nil
}

// validateRequiredIfFalse проверяет, что одно поле обязательно, если другое имеет значение false
func validateRequiredIfFalse(fl validator.FieldLevel) bool {
	param := fl.Param() // Получаем имя зависимого поля
	field := fl.Field() // Текущее поле

	// Проверяем значение зависимого поля
	otherField := fl.Parent().FieldByName(param)
	if !otherField.IsValid() {
		return false
	}

	// Если зависимое поле имеет значение false, текущее поле должно быть заполнено
	if otherField.Kind() == reflect.Bool && !otherField.Bool() {
		return field.Uint() != 0
	}

	return true
}
