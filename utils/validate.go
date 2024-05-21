package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func ValidateStruct(data interface{}) error {
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

// var validate = validator.New()

// func (v XValidator) Validate(data interface{}) []ErrorResponse {
// 	validationErrors := []ErrorResponse{}

// 	errs := validate.Struct(data)
// 	if errs != nil {
// 		for _, err := range errs.(validator.ValidationErrors) {
// 			// In this case data object is actually holding the User struct
// 			var elem ErrorResponse

// 			elem.FailedField = err.Field() // Export struct field name
// 			elem.Tag = err.Tag()           // Export struct tag
// 			elem.Value = err.Value()       // Export field value
// 			elem.Error = true

// 			validationErrors = append(validationErrors, elem)
// 		}
// 	}

// 	return validationErrors
// }
