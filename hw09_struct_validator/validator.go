package hw09structvalidator

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var ValidationErrorType = fmt.Errorf("wrong type error")
var ValidateLenError = fmt.Errorf("wrong length error")

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := strings.Builder{}
	for _, v := range v {
		res.WriteString(v.Field + ": " + v.Err.Error() + "\n")
	}
	return res.String()
}

type (
	TagsInfo map[string]string
)

func parseTagString(tagRaw string) (retInfos TagsInfo) {
	retInfos = make(TagsInfo)
	for _, tag := range strings.Split(tagRaw, " ") {
		if tag = strings.TrimSpace(tag); tag == "" {
			continue
		}
		tagParts := strings.SplitN(tag, ":", 2)
		if len(tagParts) != 2 {
			continue
		}

		tagName := strings.TrimSpace(tagParts[0])
		// if _, found := retInfos[tagName]; found {
		// 	continue
		// }
		if tagName != "validate" {
			continue
		}

		tagValuesRaw, _ := strconv.Unquote(tagParts[1])
		tagValues := make([]string, 0)
		for _, value := range strings.Split(tagValuesRaw, "|") {
			if value := strings.TrimSpace(value); value != "" {
				tagValues = append(tagValues, value)
			}
		}
		for _, tagValue := range tagValues {
			valueParts := strings.SplitN(tagValue, ":", 2)
			if len(tagParts) != 2 {
				continue
			}
			funcName := strings.TrimSpace(valueParts[0])
			funcArgs := strings.TrimSpace(valueParts[1])
			retInfos[funcName] = funcArgs
		}
	}
	return
}

func ValidateLen(fieldName string, fieldValueRaw reflect.Value, args string) ValidationError {
	result := ValidationError{
		Field: fieldName,
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(string)
	if !ok {
		result.Err = fmt.Errorf("invalid field value type")
		return result
	}
	length, err := strconv.Atoi(args)
	if err != nil {
		result.Err = fmt.Errorf("invalid value for len validation: %w", err)
		return result
	}
	if len(fieldValue) > length {
		result.Err = ValidateLenError
	}
	return result
}

var validationMap = map[string]func(fieldName string, fieldValue reflect.Value, args string) ValidationError{
	"len": ValidateLen,
}

func validationExec(fieldName string, fieldValue reflect.Value, funcName string, args string) ValidationError {
	valFunc := validationMap[funcName]
	return valFunc(fieldName, fieldValue, args)
}

func Validate(v interface{}) error {
	// Place your code here.
	var results = make(ValidationErrors, 0)

	var objType reflect.Type
	if t, ok := v.(reflect.Type); ok {
		objType = t
	} else {
		objType = reflect.ValueOf(v).Type()
	}

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objType.Kind() != reflect.Struct {
		return ValidationErrors{ValidationError{Err: ValidationErrorType}}
	}

	for fieldIdx := 0; fieldIdx < objType.NumField(); fieldIdx++ {
		field := objType.Field(fieldIdx)

		fieldValue := reflect.ValueOf(v).Field(fieldIdx)
		fmt.Println("fieldValue", fieldValue)
		tags := parseTagString(string(field.Tag))
		log.Printf("field: %v, value: %v, tags: %v\n", field.Name, fieldValue, tags)
		for funcName, args := range tags {
			valErr := validationExec(field.Name, fieldValue, funcName, args)
			log.Printf("validation result: %v\n", valErr)
			if valErr.Err != nil {
				results = append(results, valErr)
			}
		}
	}
	if len(results) > 0 {
		return results
	}
	return nil
}
