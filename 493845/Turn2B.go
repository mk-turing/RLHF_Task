package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type Parameter struct {
	Name     string
	Type     reflect.Type
	Required bool
}

type Validator interface {
	Validate(param Parameter, value string) error
}

func NewValidator(paramType reflect.Type) Validator {
	switch paramType.Kind() {
	case reflect.String:
		return &StringValidator{}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &IntValidator{}
	case reflect.Float32, reflect.Float64:
		return &FloatValidator{}
	case reflect.Bool:
		return &BoolValidator{}
	case reflect.Array, reflect.Slice:
		return &ArrayValidator{elementType: paramType.Elem()}
	case reflect.Struct:
		return &StructValidator{structType: paramType}
	default:
		return &DefaultValidator{}
	}
}

func ValidateQueryParams(params []Parameter, query url.Values) error {
	for _, param := range params {
		values := query[param.Name]
		if len(values) == 0 {
			if param.Required {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		for _, value := range values {
			validator := NewValidator(param.Type)
			if err := validator.Validate(param, value); err != nil {
				return err
			}
		}
	}
	return nil
}

type StringValidator struct{}

func (v *StringValidator) Validate(param Parameter, value string) error {
	return nil
}

type IntValidator struct{}

func (v *IntValidator) Validate(param Parameter, value string) error {
	_, err := strconv.Atoi(value)
	return err
}

type FloatValidator struct{}

func (v *FloatValidator) Validate(param Parameter, value string) error {
	_, err := strconv.ParseFloat(value, 64)
	return err
}

type BoolValidator struct{}

func (v *BoolValidator) Validate(param Parameter, value string) error {
	_, err := strconv.ParseBool(value)
	return err
}

type ArrayValidator struct {
	elementType reflect.Type
}

func (v *ArrayValidator) Validate(param Parameter, value string) error {
	var arr []interface{}
	if err := json.Unmarshal([]byte(value), &arr); err != nil {
		return fmt.Errorf("invalid array format for parameter '%s': %v", param.Name, err)
	}

	for _, element := range arr {
		elementValue := reflect.ValueOf(element)
		if elementValue.Kind() != v.elementType.Kind() {
			return fmt.Errorf("invalid element type for array parameter '%s': expected %s, got %s", param.Name, v.elementType.Kind(), elementValue.Kind())
		}
	}
	return nil
}

type StructValidator struct {
	structType reflect.Type
}

func (v *StructValidator) Validate(param Parameter, value string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		return fmt.Errorf("invalid object format for parameter '%s': %v", param.Name, err)
	}

	for i := 0; i < v.structType.NumField(); i++ {
		field := v.structType.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name