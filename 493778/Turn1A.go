package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

type ValidationRule struct {
	Name          string
	Required      bool
	DataType      string
	AllowedValues []string
}

func (rule *ValidationRule) IsValid(values url.Values) error {
	value, ok := values[rule.Name]
	if !ok && rule.Required {
		return fmt.Errorf("query parameter '%s' is required", rule.Name)
	}

	switch rule.DataType {
	case "string":
		if len(value) == 0 {
			return fmt.Errorf("query parameter '%s' cannot be empty", rule.Name)
		}
		return nil
	case "int":
		if _, err := strconv.Atoi(value[0]); err != nil {
			return fmt.Errorf("query parameter '%s' must be an integer", rule.Name)
		}
		return nil
	case "float":
		if _, err := strconv.ParseFloat(value[0], 64); err != nil {
			return fmt.Errorf("query parameter '%s' must be a float", rule.Name)
		}
		return nil
	case "bool":
		if _, err := strconv.ParseBool(value[0]); err != nil {
			return fmt.Errorf("query parameter '%s' must be a boolean", rule.Name)
		}
		return nil
	case "regex":
		pattern := rule.AllowedValues[0]
		r, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %s", pattern)
		}
		if !r.MatchString(value[0]) {
			return fmt.Errorf("query parameter '%s' does not match the regex pattern", rule.Name)
		}
		return nil
	default:
		return fmt.Errorf("invalid data type: %s", rule.DataType)
	}
}

func validateQueryParams(rules []ValidationRule) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		for _, rule := range rules {
			if err := rule.IsValid(values); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		http.Next(w, r)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	validationRules := []ValidationRule{
		{Name: "id", Required: true, DataType: "int"},
		{Name: "name", DataType: "string", AllowedValues: []string{"Alice", "Bob"}},
		{Name: "age", DataType: "int", Required: true, AllowedValues: []string{"18", "21", "25"}},
	}

	http.HandleFunc("/protected", validateQueryParams(validationRules)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Protected Resource Access Granted")
	})))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
