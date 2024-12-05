package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"os"
)

func parseAndValidate(fileName string) error {
	// Read the entire file into a byte slice
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Parse JSON and validate
	err = jsonparser.ObjectEach(fileBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "product":
			productName := ""
			productPrice := 0.0

			// Extract name
			if name, err := jsonparser.GetString(value, "name"); err == nil {
				productName = name
			} else {
				return fmt.Errorf("failed to get product name: %v", err)
			}

			// Extract price
			if price, err := jsonparser.GetFloat(value, "price"); err == nil {
				productPrice = price
			} else {
				return fmt.Errorf("failed to get product price: %v", err)
			}

			// Validate price
			if productPrice <= 0 {
				fmt.Printf("Invalid price for product: %s\n", productName)
			} else {
				fmt.Printf("Valid product: %s with price: %.2f\n", productName, productPrice)
			}

			// Extract and process categories (if any)
			_, err := jsonparser.ArrayEach(value, func(item []byte, dataType jsonparser.ValueType, offset int, err error) {
				category := string(item)
				fmt.Printf("Product: %s, Category: %s\n", productName, category)
			}, "categories")
			if err != nil && err != jsonparser.KeyPathNotFoundError {
				return fmt.Errorf("error processing categories: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error parsing file: %v", err)
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
