package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Record struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	// Add more fields as needed
}

func validateRecord(record *Record) error {
	if record.Name == "" {
		return fmt.Errorf("name is required")
	}
	if record.Age < 0 {
		return fmt.Errorf("age cannot be negative")
	}
	return nil
}

func validateFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var record Record
		err = decoder.Decode(&record)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}
		err = validateRecord(&record)
		if err != nil {
			return fmt.Errorf("validation failed for record %+v: %w", record, err)
		}
	}
	return nil
}

func main() {
	if err := validateFile("../385843/large_data_file.json"); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful")
	}
}
