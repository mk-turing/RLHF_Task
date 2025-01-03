package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Parameter defines a struct that holds the parameter name, type, required flag, and array or nested struct details
type Parameter struct {
	Name      string
	Type      reflect.Type
	Required  bool
	Array     bool
	Nested    bool
	Subtype   reflect.Type // Required for nested objects
	SubParams []Parameter  // Holds nested parameters
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
		valueStrings, ok := query[param.Name]
		if !ok {
			if param.Required {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		// Validate for an array of values
		if param.Array {
			if err := validateArrayParam(param, valueStrings); err != nil {
				return err
			}
		} else if param.Nested {
			// Validate nested parameters
			nestedValues := make(url.Values)
			for _, value := range valueStrings {
				parts := strings.Split(value, "&")
				for _, part := range parts {
					k, v, err := parseKeyValue(part)
					if err != nil {
						return fmt.Errorf("failed to parse nested parameter: %s", part)
					}
					nestedValues.Add(k, v)
				}
			}
			if err := ValidateQueryParams(param.SubParams, nestedValues); err != nil {
				return err
			}
		} else {
			// Validate for a single value
			if len(valueStrings) > 1 {
				return fmt.Errorf("parameter '%s' can have only one value", param.Name)
			}

			value := valueStrings[0]
			validator := NewValidator(param.Type)
			if err := validator.Validate(param, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateArrayParam(param Parameter, values []string) error {
	for _, value := range values {
		if err := NewValidator(param.Subtype).Validate(param, value); err != nil {
			return fmt.Errorf("invalid value for array parameter '%s': %s", param.Name, value)
		}
	}
	return nil
}

// StringValidator validates string parameters
type StringValidator struct{}

func (v *StringValidator) Validate(param Parameter, value string) error {
	if len(value) == 0 && param.Required {
		return fmt.Errorf("empty string for required parameter: %s", param.Name)
	}
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

func parseKeyValue(part string) (string, string, error) {
	kv := strings.SplitN(part, "=", 2)
	if len(kv) != 2 {
		return "", "", fmt.Errorf("invalid key-value pair: %s", part)
	}
	return kv[0], kv[1], nil
}

func main() {
	// Sample query parameters
	query := url.Values{}
	query.Add("name", "John")
	query.Add("age", "25")
	query.Add("address[street]", "123 Main St")
	query.Add("address[city]", "Metropolis")
	query.Add("products[0]", "apple")
	query.Add("products[1]", "banana")
	query.Add("isActive", "true")

	// List of expected parameters with types
	params := []Parameter{
		{Name: "name", Type: reflect.TypeOf(""), Required: true},
		{Name: "age", Type: reflect.TypeOf(0), Required: true},
		{
			Name:   "address",
			Type:   reflect.TypeOf(""),
			Array:  false,
			Nested: true,
			SubParams: []Parameter{
				{Name: "street", Type: reflect.TypeOf(""), Required: true},
				{Name: "city", Type: reflect.TypeOf(""), Required: true},
			},
		},
		{
			Name:    "products",
			Type:    reflect.TypeOf([]string{}),
			Array:   true,
			Subtype: reflect.TypeOf(""),
		},
		{Name: "isActive", Type: reflect.TypeOf(true), Required: false},
	}

	// Validate the query parameters
	if err := ValidateQueryParams(params, query); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("All parameters are valid!")
	}
}
