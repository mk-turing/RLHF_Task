package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"os"
	"strings"
)

func parseAndValidate(filename string) error {
	// Read the entire file into a byte slice
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var path string // Represent the current path in the nested structure

	// Iterate over the JSON object
	err = jsonparser.ObjectEach(fileBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// Append the key to the current path
		if path == "" {
			path = string(key)
		} else {
			path = fmt.Sprintf("%s.%s", path, string(key))
		}

		if string(key) == "user" { // Assuming objects starting with "user" are of interest
			err := validateUser(path, value)
			if err != nil {
				return err
			}
		}

		// Remove the last key from the path when an object ends
		i := strings.LastIndex(path, ".")
		if i >= 0 {
			path = path[:i]
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func validateUser(path string, value []byte) error {
	// Validate age: path like "user.data.age"
	age, err := jsonparser.GetInt(value, "age")
	if err != nil {
		return fmt.Errorf("error reading age at %s: %v", path, err)
	}
	if age <= 0 {
		return fmt.Errorf("age should be positive: %s", path)
	}
	return nil
}

func main() {
	if err := parseAndValidate("../385843/large_data_file.json"); err != nil {
		fmt.Printf("Error parsing: %s\n", err)
	} else {
		fmt.Println("File parsed and validated successfully")
	}
}
