package main

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
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
func ValidateQueryParams(params []Parameter, query url.Values, errors *[]ValidationError) error {
	for _, param := range params {
		valueStrings, ok := query[param.Name]
		if !ok {
			if param.Required {
				log.Printf("Validation Error: Missing required parameter: %s", param.Name)
				if errors != nil {
					*errors = append(*errors, ValidationError{param.Name, "Missing required parameter"})
				}
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		// Validate for an array of values
		if param.Array {
			if err := validateArrayParam(param, valueStrings, errors); err != nil {
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
						log.Printf("Validation Error: Failed to parse nested parameter: %s", part)
						if errors != nil {
							*errors = append(*errors, ValidationError{part, "Failed to parse nested parameter"})
						}
						return err
					}
					nestedValues.Add(k, v)
				}
			}
			if err := ValidateQueryParams(param.SubParams, nestedValues, errors); err != nil {
				return err
			}
		} else {
			// Validate for a single value
			if len(valueStrings) > 1 {
				log.Printf("Validation Error: Parameter '%s' can have only one value", param.Name)
				if errors != nil {
					*errors = append(*errors, ValidationError{param.Name, "Can have only one value"})
				}
				return fmt.Errorf("parameter '%s' can have only one value", param.Name)
			}

			value := valueStrings[0]
			validator := NewValidator(param.Type)
			if err := validator.Validate(param, value); err != nil {
				log.Printf("Validation Error: Invalid value for parameter '%s': %s", param.Name, value)
				if errors != nil {
					*errors = append(*errors, ValidationError{param.Name, err.Error()})
				}
				return err
			}
		}
	}
	return nil
}

func validateArrayParam(param Parameter, values []string, errors *[]ValidationError) error {
	for i, value := range values {
		if err := NewValidator(param.Subtype).Validate(param, value); err != nil {
			log.Printf("Validation Error: Invalid value for array parameter '%s[%d]': %s", param.Name, i, value)
			if errors != nil {
				*errors = append(*errors, ValidationError{fmt.Sprintf("%s[%d]", param.Name, i), err.Error()})
			}
			return err
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

// ValidationError represents an error encountered during validation
type ValidationError struct {
	Field   string
	Message string
}

// GenerateReport generates a detailed report of all validation errors
func GenerateReport(errors []ValidationError) {
	if len(errors) == 0 {
		fmt.Println("No validation errors found.")
		return
	}

	fmt.Println("Validation Report:")
	for _, error := range errors {
		fmt.Printf("- Field: %s, Message: %s\n", error.Field, error.Message)
	}
}

// ProcessQueryParametersConcurrently validates query parameters concurrently
func ProcessQueryParametersConcurrently(params []Parameter, queries []url.Values, wg *sync.WaitGroup, errors *[]ValidationError) {
	defer wg.Done()

	for _, query := range queries {
		var queryErrors []ValidationError
		if err := ValidateQueryParams(params, query, &queryErrors); err != nil {
			log.Printf("Validation failed for query: %+v", query)
			if errors != nil {
				*errors = append(*errors, queryErrors...)
			}
		} else {
			log.Printf("All parameters are valid for query: %+v", query)
		}
	}
}

func main() {
	// Sample query parameters
	queries := []url.Values{
		{
			"name":        {"John"},
			"age":         {"25"},
			"address":     {"123 Main St", "Metropolis"},
			"products":    {"apple", "banana"},
			"isActive":    {"true"},
			"notAInteger": {"notAnIntegerValue"},
		},
		{
			"name":     {"Alice"},
			"age":      {"30"},
			"address":  {"456 Elm St", "Springfield"},
			"products": {"orange", "grape"},
			"isActive": {"false"},
		},
		// Add more queries as needed
	}

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
		{Name: "notAInteger", Type: reflect.TypeOf(0), Required: false},
	}

	// Start a wait group to manage concurrent processing
	var wg sync.WaitGroup
	var validationErrors []ValidationError

	// Process queries concurrently
	start := time.Now()
	for _, query := range queries {
		wg.Add(1)
		go ProcessQueryParametersConcurrently(params, []url.Values{query}, &wg, &validationErrors)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	elapsed := time.Since(start)

	// Generate and print the report
	GenerateReport(validationErrors)
	fmt.Printf("Validation completed in %s\n", elapsed)
}
