package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type queryValidator struct {
	params map[string][]validator
}

type validator interface {
	Validate(string) error
}

type requiredValidator struct{}

func (v requiredValidator) Validate(s string) error {
	if s == "" {
		return fmt.Errorf("parameter is required")
	}
	return nil
}

type intValidator struct{}

func (v intValidator) Validate(s string) error {
	_, err := strconv.Atoi(s)
	return err
}

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

func newQueryValidator() *queryValidator {
	return &queryValidator{params: make(map[string][]validator)}
}

func (v *queryValidator) AddRule(param string, validators ...validator) {
	v.params[param] = append(v.params[param], validators...)
}

func (v *queryValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for param, validators := range v.params {
			values, ok := r.URL.Query()[param]
			if !ok || len(values) == 0 {
				continue
			}

			for _, value := range values {
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
	v.AddRule("limit", requiredValidator{}, intValidator{})
	v.AddRule("order", allowedValuesValidator{allowed: []string{"asc", "desc"}})

	http.Handle("/api/data", v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Data endpoint reached")
	})))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
