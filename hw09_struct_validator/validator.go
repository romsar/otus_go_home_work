package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const ValidationTag = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotValidatable         = errors.New("that entity cannot be validated")
	ErrInvalidValidationRule  = errors.New("invalid validation rule")
	ErrInvalidValidationValue = errors.New("invalid value")
)

func (ve ValidationErrors) Error() string {
	var msg strings.Builder

	for _, v := range ve {
		if msg.Len() != 0 {
			msg.WriteString(", ")
		}
		msg.WriteString(v.Error())
	}

	return msg.String()
}

func (ve *ValidationErrors) Add(field string, err error) {
	*ve = append(*ve, ValidationError{
		Field: field,
		Err:   err,
	})
}

func (v ValidationError) Error() string {
	if v.Err == nil {
		return ""
	}
	return v.Err.Error()
}

func Validate(v interface{}) error {
	refVal := reflect.ValueOf(v)

	if checkStruct(refVal) {
		return ErrNotValidatable
	}

	ve := make(ValidationErrors, 0)

	for i := 0; i < refVal.Type().NumField(); i++ {
		typeField := refVal.Type().Field(i)

		// if private field = skip
		if !typeField.IsExported() {
			continue
		}

		// if no validation rules = skip
		vRulesStr := validationTagValue(typeField)
		if vRulesStr == "" {
			continue
		}

		field := refVal.Field(i)
		vRules := validationRules(vRulesStr)

		if err := iterateAndValidate(vRules, typeField, field, &ve); err != nil {
			return err
		}
	}

	if !hasValidationErrors(ve) {
		return nil
	}

	return ve
}

func iterateAndValidate(
	vRules []string,
	typeField reflect.StructField,
	field reflect.Value,
	ve *ValidationErrors,
) error {
	for _, r := range vRules {
		if isIterable(typeField.Type.Kind()) {
			if err := validateStruct(field, r, ve, typeField); err != nil {
				return err
			}
		} else {
			if err := validateSingle(typeField, field, r, ve); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSingle(typeField reflect.StructField, field reflect.Value, r string, ve *ValidationErrors) error {
	if err := doValidate(typeField.Type.Kind(), field, r); err != nil {
		if errors.Is(err, ErrInvalidValidationValue) {
			ve.Add(typeField.Name, err)
			return nil
		}
		return err
	}
	return nil
}

func validateStruct(field reflect.Value, r string, ve *ValidationErrors, typeField reflect.StructField) error {
	sliceVal := reflect.ValueOf(field.Interface())
	for i := 0; i < sliceVal.Len(); i++ {
		sliceElem := sliceVal.Index(i)
		err := doValidate(sliceElem.Kind(), sliceElem, r)
		if err == nil {
			continue
		}
		if errors.Is(err, ErrInvalidValidationValue) {
			ve.Add(typeField.Name, err)
			continue
		}
		return err
	}
	return nil
}

func checkStruct(t reflect.Value) bool {
	return t.Type().Kind() != reflect.Struct
}

func validationTagValue(typeField reflect.StructField) string {
	return typeField.Tag.Get(ValidationTag)
}

func validationRules(rules string) []string {
	return strings.Split(rules, "|")
}

func hasValidationErrors(ve ValidationErrors) bool {
	return len(ve) > 0
}

func isIterable(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

func doValidate(kind reflect.Kind, field reflect.Value, r string) error {
	ruleName, ruleValStr, ok := splitRule(r)
	if !ok {
		return ErrInvalidValidationRule
	}

	err := doValidateValue(kind, field, ruleName, ruleValStr)
	if err != nil {
		return err
	}

	return nil
}

func doValidateValue(kind reflect.Kind, field reflect.Value, ruleName string, ruleValStr string) error {
	if kind == reflect.String {
		fieldValStr := field.String()
		switch ruleName {
		case "len":
			return validateStringLen(ruleValStr, fieldValStr)
		case "in":
			return validateStringIn(ruleValStr, fieldValStr)
		case "regexp":
			return validateStringRegexp(ruleValStr, fieldValStr)
		}
	} else if kind == reflect.Int {
		fieldValInt := field.Int()
		switch ruleName {
		case "min":
			return validateIntMin(ruleValStr, fieldValInt)
		case "max":
			return validateIntMax(ruleValStr, fieldValInt)
		case "in":
			return validateIntIn(ruleValStr, fieldValInt)
		}
	}

	return nil
}

func validateIntIn(inStr string, val int64) error {
	in := strings.Split(inStr, ",")
	inInt := make([]int64, 0, len(in))

	for _, v := range in {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return ErrInvalidValidationRule
		}
		inInt = append(inInt, int64(vInt))
	}

	var ok bool
	for _, n := range inInt {
		if n == val {
			ok = true
			break
		}
	}
	if ok {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value (%d) is not in (%s)",
		val,
		inStr,
	)
	return wrapValidationErr(valErrMsg)
}

func validateIntMax(maxStr string, val int64) error {
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return ErrInvalidValidationRule
	}
	if int64(max) >= val {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value (%d) is more than (%d)",
		val,
		max,
	)
	return wrapValidationErr(valErrMsg)
}

func validateIntMin(minStr string, val int64) error {
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return ErrInvalidValidationRule
	}
	if int64(min) <= val {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value (%d) is less than (%d)",
		val,
		min,
	)
	return wrapValidationErr(valErrMsg)
}

func validateStringRegexp(r string, val string) error {
	matched, err := regexp.MatchString(r, val)
	if matched && err == nil {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value (%s) is not matched regexp (%s)",
		val,
		r,
	)
	return wrapValidationErr(valErrMsg)
}

func validateStringIn(inStr string, val string) error {
	in := strings.Split(inStr, ",")

	var ok bool
	for _, n := range in {
		if n == val {
			ok = true
			break
		}
	}
	if ok {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value (%s) is not in (%s)",
		val,
		inStr,
	)
	return wrapValidationErr(valErrMsg)
}

func validateStringLen(ruleValStr string, fieldValStr string) error {
	ruleVal, err := strconv.Atoi(ruleValStr)
	if err != nil {
		return ErrInvalidValidationRule
	}

	if ruleVal == len(fieldValStr) {
		return nil
	}

	valErrMsg := fmt.Sprintf(
		"value length (%d) is not match required length (%s)",
		len(fieldValStr),
		ruleValStr,
	)
	return wrapValidationErr(valErrMsg)
}

func splitRule(r string) (string, string, bool) {
	splitted := strings.Split(r, ":")
	if len(splitted) != 2 {
		return "", "", false
	}
	ruleName, ruleValStr := splitted[0], splitted[1]
	return ruleName, ruleValStr, true
}

func wrapValidationErr(str string) error {
	return fmt.Errorf("%w: %s", ErrInvalidValidationValue, str)
}
