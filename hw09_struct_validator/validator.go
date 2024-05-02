package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string

	Err error
}

var ErrValidate = fmt.Errorf("wrong type error")

var ErrValidateLen = fmt.Errorf("wrong length error")

var ErrValidateRegexp = fmt.Errorf("string does not match regexp")

var ErrValidateMin = fmt.Errorf("value is less than min")

var ErrValidateMax = fmt.Errorf("value is more than max")

var ErrValidateIn = fmt.Errorf("field not in set")

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {

	res := strings.Builder{}

	var err error

	for _, v := range v {

		err = fmt.Errorf("%s: %w", v.Field, v.Err)

		res.WriteString(err.Error())

	}

	return err.Error()

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

		Err: nil,
	}

	var fieldValues []string

	switch fieldValueRaw.Kind() { //nolint:exhaustive

	case reflect.String:

		fieldValues = []string{fieldValueRaw.String()}

	case reflect.Slice:

		if fieldValueRaw.Type() == reflect.TypeOf(fieldValues) {

			fieldValues = make([]string, 0)

			fieldValues = append(fieldValues, fieldValueRaw.Interface().([]string)...)

		} else {

			result.Err = fmt.Errorf("invalid field value type")

			return result

		}

	default:

		result.Err = fmt.Errorf("invalid field value type")

		return result

	}

	for _, fieldValue := range fieldValues {

		length, err := strconv.Atoi(args)

		if err != nil {

			result.Err = fmt.Errorf("invalid value for len validation: %w", err)

			return result

		}

		if len(fieldValue) != length {

			result.Err = ErrValidateLen

		}

	}

	return result

}

func ValidateRegexp(fieldName string, fieldValueRaw reflect.Value, reStr string) ValidationError {

	result := ValidationError{

		Field: fieldName,

		Err: nil,
	}

	var fieldValues []string

	switch fieldValueRaw.Kind() { //nolint:exhaustive

	case reflect.String:

		fieldValues = []string{fieldValueRaw.String()}

	case reflect.Slice:

		if fieldValueRaw.Type() == reflect.TypeOf(fieldValues) {

			fieldValues = make([]string, 0)

			fieldValues = append(fieldValues, fieldValueRaw.Interface().([]string)...)

		} else {

			result.Err = fmt.Errorf("invalid field value type")

			return result

		}

	default:

		result.Err = fmt.Errorf("invalid field value type")

		return result

	}

	re, err := regexp.Compile(reStr)

	if err != nil {

		result.Err = fmt.Errorf("invalid regexp: %w", err)

		return result

	}

	for _, fieldValue := range fieldValues {

		if !re.MatchString(fieldValue) {

			result.Err = ErrValidateRegexp

			return result

		}

	}

	return result

}

var validationMap = map[string]func(fieldName string, fieldValue reflect.Value, args string) ValidationError{

	"len": ValidateLen,

	"regexp": ValidateRegexp,

	"in": ValidateIn,

	"min": ValidateMin,

	"max": ValidateMax,
}

func ValidateIn(fieldName string, fieldValueRaw reflect.Value, args string) ValidationError {

	result := ValidationError{

		Field: fieldName,

		Err: nil,
	}

	var fieldValues []any

	switch fieldValueRaw.Kind() { //nolint:exhaustive

	case reflect.String:

		fieldValues = []any{fieldValueRaw.String()}

	case reflect.Int:

		fieldValues = []any{fieldValueRaw.Int()}

	case reflect.Slice:

		for i := 0; i < fieldValueRaw.Len(); i++ {

			fieldValues = append(fieldValues, fieldValueRaw.Index(i).Interface())

		}

	default:

		result.Err = fmt.Errorf("invalid field value type")

		return result

	}

	for _, fieldValue := range fieldValues {

		fieldValueStr := fmt.Sprintf("%v", fieldValue)

		for _, arg := range strings.Split(args, ",") {

			if arg == fieldValueStr {

				return result

			}

		}

	}

	result.Err = ErrValidateIn

	return result

}

func ValidateMin(fieldName string, fieldValueRaw reflect.Value, args string) ValidationError {

	result := ValidationError{

		Field: fieldName,

		Err: nil,
	}

	var fieldValues []int

	switch fieldValueRaw.Kind() { //nolint:exhaustive

	case reflect.Int:

		fieldValues = []int{int(fieldValueRaw.Int())}

	case reflect.Slice:

		if fieldValueRaw.Type() == reflect.TypeOf(fieldValues) {

			fieldValues = make([]int, 0)

			fieldValues = append(fieldValues, fieldValueRaw.Interface().([]int)...)

		} else {

			result.Err = fmt.Errorf("invalid field value type")

			return result

		}

	}

	min, err := strconv.Atoi(args)

	if err != nil {

		result.Err = fmt.Errorf("invalid value for min validation: %w", err)

		return result

	}

	for _, fieldValue := range fieldValues {

		if fieldValue < min {

			result.Err = ErrValidateMin

		}

	}

	return result

}

func ValidateMax(fieldName string, fieldValueRaw reflect.Value, args string) ValidationError {

	result := ValidationError{

		Field: fieldName,

		Err: nil,
	}

	var fieldValues []int

	switch fieldValueRaw.Kind() { //nolint:exhaustive

	case reflect.Int:

		fieldValues = []int{int(fieldValueRaw.Int())}

	case reflect.Slice:

		if fieldValueRaw.Type() == reflect.TypeOf(fieldValues) {

			fieldValues = make([]int, 0)

			fieldValues = append(fieldValues, fieldValueRaw.Interface().([]int)...)

		} else {

			result.Err = fmt.Errorf("invalid field value type")

			return result

		}

	}

	max, err := strconv.Atoi(args)

	if err != nil {

		result.Err = fmt.Errorf("invalid value for max validation: %w", err)

		return result

	}

	for _, fieldValue := range fieldValues {

		if fieldValue > max {

			result.Err = ErrValidateMax

		}

	}

	return result

}

func validationExec(fieldName string, fieldValue reflect.Value, funcName string, args string) ValidationError {

	valFunc, ok := validationMap[funcName]

	if !ok {

		return ValidationError{Err: errors.New("unknown validation function")}

	}

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

		return ValidationErrors{ValidationError{Err: ErrValidate}}

	}

	for fieldIdx := 0; fieldIdx < objType.NumField(); fieldIdx++ {

		field := objType.Field(fieldIdx)

		fieldValue := reflect.ValueOf(v).Field(fieldIdx)

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
