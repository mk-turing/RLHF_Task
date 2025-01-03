package main

import (
	"fmt"
)

// QueryParam is a struct representing the request URL's query parameters
type QueryParam struct {
	Key   string
	Value string
}

// ValidateQueryParams checks that query parameters adhere to specified constraints
func ValidateQueryParams(url string, params []QueryParam) ([]error, bool) {
	parsedURL, err := url.Parse(url)
	if err != nil {
		return []error{err}, false
	}

	query := parsedURL.Query()
	errors := make([]error, len(params))
	success := true

	for i, param := range params {
		value, _ := query.Get(param.Key)
		if value != param.Value {
			errors[i] = fmt.Errorf("expected %q for key %q, got %q", param.Value, param.Key, value)
			success = false
		}
	}

	return errors, success
}

func main() {
	url := "https://example.com/page?id=123&name=johndoe"
	params := []QueryParam{
		{"id", "123"},
		{"name", "johndoe"},
		{"flag", "false"}, // Example parameter that doesn't exist
	}

	errors, success := ValidateQueryParams(url, params)
	if success {
		fmt.Println("All validation passed.")
	} else {
		fmt.Println("Validation failed for:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
