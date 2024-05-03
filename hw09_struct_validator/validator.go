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

var (
	ErrValidate       = errors.New("wrong type error")
	ErrValidateLen    = errors.New("wrong length error")
	ErrValidateRegexp = errors.New("string does not match regexp")
	ErrValidateMin    = errors.New("value is less than min")
	ErrValidateMax    = errors.New("value is more than max")
	ErrValidateIn     = errors.New("field not in set")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := strings.Builder{}
	for _, v := range v {
		res.WriteString(fmt.Errorf("%s: %w", v.Field, v.Err).Error() + "\n")
	}
	return res.String()
}

func (v ValidationErrors) Unwrap() []error {
	res := []error{}
	for _, v := range v {
		res = append(res, v.Err)
	}
	return res
}

type (
	TagsInfo map[string]string
)

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
		Err:   nil,
	}
	fieldValue, ok := fieldValueRaw.Interface().(string)
	if !ok {
		result.Err = ErrValidate
		return result
	}
	length, err := strconv.Atoi(args)
	if err != nil {
		result.Err = fmt.Errorf("invalid value for len validation: %w", err)
		return result
	}
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
	}
	return result
}

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
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		tags := parseTagString(validateTag)
		for funcName, args := range tags {
			if fieldValue.Kind() == reflect.Slice {
				for i := 0; i < fieldValue.Len(); i++ {
					valErr := validationExec(field, fieldValue.Index(i), funcName, args)
					if valErr.Err != nil {
						results = append(results, valErr)
					}
					if errors.Is(valErr.Err, ErrValidate) {
						break
					}
				}
			} else {
				valErr := validationExec(field, fieldValue, funcName, args)
				if valErr.Err != nil {
					results = append(results, valErr)
				}
			}
		}
	}
	if len(results) > 0 {
		return results
	}
	return nil
}
