package main
import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Parameter struct {
	Name     string
	Type     reflect.Type
	Required bool
	Array    bool
	Nested   bool
	Subtype   reflect.Type
	SubParams []Parameter
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
	default:
		return &DefaultValidator{}
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func ValidateQueryParamsConcurrently(params []Parameter, query url.Values, concurrency int, errors chan<- ValidationError) {
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for param, _ := range query {
				if err := ValidateQueryParam(params, param, query[param], errors); err != nil {
					log.Printf("Validation Error: %s", err.Error())
				}
			}
		}()
	}

	wg.Wait()
	close(errors)
}

func ValidateQueryParam(params []Parameter, param string, values []string, errors chan<- ValidationError) error {
	for _, p := range params {
		if p.Name == param {
			if err := validateParam(p, values); err != nil {
				errors <- ValidationError{param, err.Error()}
				return err
			}
			return nil
		}
	}
	errors <- ValidationError{param, "Parameter not found"}
	return fmt.Errorf("parameter '%s' not found", param)
}

func validateParam(param Parameter, values []string) error {
	if len(values) == 0 {
		if param.Required {
			return fmt.Errorf("missing required parameter: %s", param.Name)
		}
		return nil
	}

	if param.Array {
		for _, value := range values {
			if err := NewValidator(param.Subtype).Validate(param, value); err != nil {
				return err
			}
		}
		return nil
	}

	if param.Nested {
		nestedValues := make(url.Values)
		for _, value := range values {
			parts := strings.Split(value, "&")
			for _, part := range parts {
				k, v, err := parseKeyValue(part)
				if err != nil {
					return err
				}
				nestedValues.Add(k, v)
			}
		}
		return ValidateQueryParam(param.SubParams, "", nestedValues[""], nil)
	}

	if len(values) > 1 {
		return fmt.Errorf("parameter '%s' can have only one value", param.Name)
	}

	return NewValidator(param.Type).Validate(param, values[0])
}

func parseKeyValue(part string) (string, string, error) {
	kv := strings.SplitN(part, "=", 2)
	if len(kv) != 2 {
		return "", "", fmt.Errorf("invalid key-value pair: %s", part)
	}