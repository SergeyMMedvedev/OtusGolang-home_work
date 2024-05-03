package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

<<<<<<< HEAD
var ValidationErrorType = fmt.Errorf("wrong type error")
var ValidateLenError = fmt.Errorf("wrong length error")
=======
var (
	ErrValidate       = errors.New("wrong type error")
	ErrValidateLen    = errors.New("wrong length error")
	ErrValidateRegexp = errors.New("string does not match regexp")
	ErrValidateMin    = errors.New("value is less than min")
	ErrValidateMax    = errors.New("value is more than max")
	ErrValidateIn     = errors.New("field not in set")
)
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := strings.Builder{}
	for _, v := range v {
<<<<<<< HEAD
		res.WriteString(v.Field + ": " + v.Err.Error() + "\n")
=======
		res.WriteString(fmt.Errorf("%s: %w", v.Field, v.Err).Error() + "\n")
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
	}
	return res.String()
}

<<<<<<< HEAD
=======
func (v ValidationErrors) Unwrap() []error {
	res := []error{}
	for _, v := range v {
		res = append(res, v.Err)
	}
	return res
}

>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
type (
	TagsInfo map[string]string
)

<<<<<<< HEAD
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
=======
func parseTagString(tag string) (retInfos TagsInfo) {
	retInfos = make(TagsInfo)
	tagValues := make([]string, 0)
	for _, value := range strings.Split(tag, "|") {
		if value := strings.TrimSpace(value); value != "" {
			tagValues = append(tagValues, value)
		}
	}
	for _, tagValue := range tagValues {
		valueParts := strings.SplitN(tagValue, ":", 2)
		funcName := strings.TrimSpace(valueParts[0])
		funcArgs := strings.TrimSpace(valueParts[1])
		retInfos[funcName] = funcArgs
	}
	return retInfos
}

func ValidateLen(field reflect.StructField, fieldValueRaw reflect.Value, args string) ValidationError {
	result := ValidationError{
		Field: field.Name,
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(string)
	if !ok {
<<<<<<< HEAD
		result.Err = fmt.Errorf("invalid field value type")
=======
		result.Err = ErrValidate
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
		return result
	}
	length, err := strconv.Atoi(args)
	if err != nil {
		result.Err = fmt.Errorf("invalid value for len validation: %w", err)
		return result
	}
<<<<<<< HEAD
	if len(fieldValue) > length {
		result.Err = ValidateLenError
=======
	if len(fieldValue) != length {
		result.Err = ErrValidateLen
	}

	return result
}

func ValidateRegexp(field reflect.StructField, fieldValueRaw reflect.Value, reStr string) ValidationError {
	result := ValidationError{
		Field: field.Name,
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(string)
	if !ok {
		result.Err = ErrValidate
		return result
	}
	re, err := regexp.Compile(reStr)
	if err != nil {
		result.Err = fmt.Errorf("invalid regexp: %w", err)
		return result
	}
	if !re.MatchString(fieldValue) {
		result.Err = ErrValidateRegexp
		return result
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
	}
	return result
}

<<<<<<< HEAD
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
=======
var validationMap = map[string]func(field reflect.StructField, fieldValue reflect.Value, args string) ValidationError{
	"len":    ValidateLen,
	"regexp": ValidateRegexp,
	"in":     ValidateIn,
	"min":    ValidateMin,
	"max":    ValidateMax,
}

func ValidateIn(field reflect.StructField, fieldValueRaw reflect.Value, args string) ValidationError {
	result := ValidationError{
		Field: field.Name,
		Err:   nil,
	}
	fieldValueStr := fmt.Sprintf("%v", fieldValueRaw)
	for _, arg := range strings.Split(args, ",") {
		if arg == fieldValueStr {
			return result
		}
	}
	result.Err = ErrValidateIn
	return result
}

func ValidateMin(field reflect.StructField, fieldValueRaw reflect.Value, args string) ValidationError {
	result := ValidationError{
		Field: field.Name,
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(int)
	if !ok {
		result.Err = ErrValidate
		return result
	}
	min, err := strconv.Atoi(args)
	if err != nil {
		result.Err = fmt.Errorf("invalid value for min validation: %w", err)
		return result
	}
	if fieldValue < min {
		result.Err = ErrValidateMin
	}
	return result
}

func ValidateMax(field reflect.StructField, fieldValueRaw reflect.Value, args string) ValidationError {
	result := ValidationError{
		Field: field.Name,
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(int)
	if !ok {
		result.Err = ErrValidate
		return result
	}
	max, err := strconv.Atoi(args)
	if err != nil {
		result.Err = fmt.Errorf("invalid value for max validation: %w", err)
		return result
	}
	if fieldValue > max {
		result.Err = ErrValidateMax
	}
	return result
}

func validationExec(field reflect.StructField, fieldValue reflect.Value, funcName string, args string) ValidationError {
	valFunc, ok := validationMap[funcName]
	if !ok {
		return ValidationError{Err: errors.New("unknown validation function")}
	}
	return valFunc(field, fieldValue, args)
}

func Validate(v interface{}) error {
	results := make(ValidationErrors, 0)
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac

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
<<<<<<< HEAD
		return ValidationErrors{ValidationError{Err: ValidationErrorType}}
=======
		return ValidationErrors{ValidationError{Err: ErrValidate}}
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
	}

	for fieldIdx := 0; fieldIdx < objType.NumField(); fieldIdx++ {
		field := objType.Field(fieldIdx)
<<<<<<< HEAD

		fieldValue := reflect.ValueOf(v).Field(fieldIdx)
		fmt.Println("fieldValue", fieldValue)
		tags := parseTagString(string(field.Tag))
		log.Printf("field: %v, value: %v, tags: %v\n", field.Name, fieldValue, tags)
		for funcName, args := range tags {
			valErr := validationExec(field.Name, fieldValue, funcName, args)
			log.Printf("validation result: %v\n", valErr)
			if valErr.Err != nil {
				results = append(results, valErr)
=======
		fieldValue := reflect.ValueOf(v).Field(fieldIdx)
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		tags := parseTagString(validateTag)
		for funcName, args := range tags {
			switch fieldValue.Kind() {
			case reflect.Slice:
				for i := 0; i < fieldValue.Len(); i++ {
					valErr := validationExec(field, fieldValue.Index(i), funcName, args)
					if valErr.Err != nil {
						results = append(results, valErr)
					}
					if errors.Is(valErr.Err, ErrValidate) {
						break
					}
				}
			default:
				valErr := validationExec(field, fieldValue, funcName, args)
				if valErr.Err != nil {
					results = append(results, valErr)
				}
>>>>>>> 4ef8f87dbbc52baf020a477c9eee93366a1938ac
			}
		}
	}
	if len(results) > 0 {
		return results
	}
	return nil
}
