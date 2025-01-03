package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// Parameter defines a struct that holds the parameter name, type, and required flag
type Parameter struct {
	Name     string
	Type     reflect.Type
	Required bool
}

// Validator defines an interface that validates a URL query parameter
type Validator interface {
	Validate(param Parameter, value string) error
}

// NewValidator creates a new Validator based on the parameter type
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
	default:
		return &DefaultValidator{}
	}
}

// ValidateQueryParams validates the URL query parameters against the given parameter list
func ValidateQueryParams(params []Parameter, query url.Values) error {
	for _, param := range params {
		value, ok := query[param.Name]
		if !ok {
			if param.Required {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		if len(value) > 1 {
			return fmt.Errorf("parameter '%s' can have only one value", param.Name)
		}

		validator := NewValidator(param.Type)
		if err := validator.Validate(param, value[0]); err != nil {
			return err
		}
	}
	return nil
}

// StringValidator validates string parameters
type StringValidator struct{}

func (v *StringValidator) Validate(param Parameter, value string) error {
	return nil
}

// IntValidator validates integer parameters
type IntValidator struct{}

func (v *IntValidator) Validate(param Parameter, value string) error {
	_, err := strconv.Atoi(value)
	return err
}

// FloatValidator validates float parameters
type FloatValidator struct{}

func (v *FloatValidator) Validate(param Parameter, value string) error {
	_, err := strconv.ParseFloat(value, 64)
	return err
}

// BoolValidator validates boolean parameters
type BoolValidator struct{}

func (v *BoolValidator) Validate(param Parameter, value string) error {
	_, err := strconv.ParseBool(value)
	return err
}

// DefaultValidator validates parameters of unrecognized types
type DefaultValidator struct{}

func (v *DefaultValidator) Validate(param Parameter, value string) error {
	return fmt.Errorf("unrecognized parameter type: %s", param.Type)
}

func main() {
	// Sample query parameters
	query := url.Values{}
	query.Add("name", "John")
	query.Add("age", "25")
	query.Add("height", "5.9")
	query.Add("isActive", "true")

	// List of expected parameters with types
	params := []Parameter{
		{Name: "name", Type: reflect.TypeOf(""), Required: true},
		{Name: "age", Type: reflect.TypeOf(0), Required: true},
		{Name: "height", Type: reflect.TypeOf(0.0), Required: false},
		{Name: "isActive", Type: reflect.TypeOf(true), Required: false},
	}

	// Validate the query parameters
	if err := ValidateQueryParams(params, query); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("All parameters are valid!")
	}
}
