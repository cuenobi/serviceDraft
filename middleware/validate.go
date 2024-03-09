package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s interface{}) error {
	return validate.Struct(s)
}

func ErrorResponse(err error, targetStruct interface{}) map[string]interface{} {
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]interface{})

	targetType := reflect.TypeOf(targetStruct)

	for _, fieldError := range validationErrors {
		fieldName := fieldError.Field()
		jsonTag := findJSONTag(targetType, fieldName)
		errorMessage := fmt.Sprintf("%s is %s", jsonTag, fieldError.Tag())
		errorMessages[jsonTag] = errorMessage
	}

	return errorMessages
}

func findJSONTag(t reflect.Type, fieldName string) string {
	field, _ := t.FieldByName(fieldName)
	tag := field.Tag.Get("json")
	if tag != "" && tag != "-" {
		tagParts := strings.Split(tag, ",")
		return tagParts[0]
	}
	return ""
}
