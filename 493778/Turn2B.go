
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type queryValidator struct {
	params map[string]func(*http.Request) []validator
}

type validator interface {
	Validate(string) error
}

// Required validator remains the same
type requiredValidator struct{}

func (v requiredValidator) Validate(s string) error {
	if s == "" {
		return fmt.Errorf("parameter is required")
	}
	return nil
}

// Int validator remains the same
type intValidator struct{}

func (v intValidator) Validate(s string) error {
	_, err := strconv.Atoi(s)
	return err
}

// AllowedValuesValidator remains the same
type allowedValuesValidator struct {
	allowed []string
}

func (v allowedValuesValidator) Validate(s string) error {
	for _, a := range v.allowed {
		if s == a {
			return nil
		}
	}
	return fmt.Errorf("parameter value is not allowed")
}

// New function to create a conditional validator
func ConditionalValidator(condition func(*http.Request) bool, validators ...validator) func(*http.Request) []validator {
	return func(r *http.Request) []validator {
		if condition(r) {
			return validators
		}
		return nil
	}
}

func newQueryValidator() *queryValidator {
	return &queryValidator{params: make(map[string]func(*http.Request) []validator)}
}

func (v *queryValidator) AddRule(param string, rule func(*http.Request) []validator) {
	v.params[param] = rule
}

func (v *queryValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for param, rule := range v.params {
			values, ok := r.URL.Query()[param]
			if !ok || len(values) == 0 {
				continue
			}

			for _, value := range values {
				validators := rule(r)
				if validators == nil {
					continue // Skip validation if the condition doesn't match
				}
				for _, validator := range validators {
					if err := validator.Validate(value); err != nil {
						http.Error(w, fmt.Sprintf("Invalid query parameter %s: %v", param, err), http.StatusBadRequest)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	v := newQueryValidator()

	// Example of conditional validation: 'limit' is required when 'page' is present
	v.AddRule("limit", ConditionalValidator(
		func(r *http.Request) bool {
			_, ok := r.URL.Query()["page"]
			return ok
		},
		requiredValidator{}, intValidator{},
	))

	// Example of conditional validation: 'order' is required when 'sort_by' is present
	v.AddRule("order", ConditionalValidator(
		func(r *http.Request) bool {
			_, ok := r.URL.Query()["sort_by"]
			return ok
		},
		requiredValidator{}, allowedValuesValidator{allowed: []string{"asc", "desc"}},
	))

	// Example of simple validation: 'sort_by' can be 'name' or 'age'
	v.AddRule("sort_by", allowedValuesValidator{allowed: []string{"name", "age"}})

	http.Handle("/api/data", v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Data endpoint reached")
	})))