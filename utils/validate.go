package utils

import (
	"fmt"
	"reflect"
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
	return nil
}
