package errors

import (
	"fmt"
	"reflect"
)

// New accept a string as parameter and return an error
func New(message string) error {
	return fmt.Errorf(message)
}

// FieldsVerifier accept a struct pointer and verifier required fields
// An error is returned when field is an empty string or nil
// It does not verifier integer or float
// omitfields specify the field name that will be omitted
func FieldsVerifier(structPointer interface{}, omitfields ...string) error {
	elements := reflect.ValueOf(structPointer).Elem()

	for i := 0; i < elements.NumField(); i++ {
		fieldName := elements.Type().Field(i).Name
		if omit(fieldName, omitfields) {
			continue
		}
		field := elements.Field(i)
		fieldType := field.Type()

		switch field.Kind() {
		case reflect.Struct:
			if field.NumField() > 0 {
				if field.IsZero() {
					return errorMessage(fieldName, fieldType.String())
				}
			}
		case reflect.Map, reflect.Array, reflect.Slice:
			if field.Len() == 0 {
				return errorMessage(fieldName, fieldType.String())
			}

		case reflect.Ptr:
			if field.IsNil() {
				return errorMessage(fieldName, fieldType.String())
			}
		case reflect.Int, reflect.Float32, reflect.Float64:
			break
		default:
			if field.IsZero() {
				return errorMessage(fieldName, fieldType.String())
			}
		}
	}
	return nil
}

// omit Loop over the omitfieds to see if the field is not required
func omit(fieldName string, omitfields []string) bool {
	for _, omitfield := range omitfields {
		if fieldName == omitfield {
			return true
		}
	}
	return false
}

func errorMessage(fieldName string, fieldType string) error {
	return fmt.Errorf("required field %s type of %s is empty", fieldName, fieldType)
}
