package main

import (
	"errors"
	"fmt"
)

type Data map[string]interface{}
type RollbackMap map[string]interface{}

// Define custom error types
type ValidationError struct {
	Key string
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for key '%s': %v", e.Key, e.Error)
}

type TransformationError struct {
	Key string
	Err error
}

func (e *TransformationError) Error() string {
	return fmt.Sprintf("Transformation error for key '%s': %v", e.Key, e.Error)
}

// Example transformation step
func transformAdd(data Data, rollbackMap *RollbackMap, key string, value interface{}) error {
	if data[key] != nil {
		return &ValidationError{Key: key, Err: fmt.Errorf("Key %s already exists", key)}
	}
	data[key] = value
	(*rollbackMap)[key] = nil
	return nil
}

// Example transformation step
func transformMultiply(data Data, rollbackMap *RollbackMap, key string) error {
	value, ok := data[key].(int)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for multiplication")}
	}
	data[key] = value * 2
	(*rollbackMap)[key] = value
	return nil
}

// Updated rollback function with error handling
func rollback(data Data, rollbackMap RollbackMap) {
	for key, oldValue := range rollbackMap {
		switch v := oldValue.(type) {
		case int:
			data[key] = v
		default:
			delete(data, key)
		}
	}
	rollbackMap = nil
}

func main() {
	data := Data{
		"a": 10,
		"b": "hello",
	}
	rollbackMap := RollbackMap{}

	// Transform the data
	err := transformAdd(data, &rollbackMap, "c", 20)
	if err != nil {
		handleError(data, rollbackMap, err)
		return
	}

	err = transformMultiply(data, &rollbackMap, "a")
	if err != nil {
		handleError(data, rollbackMap, err)
		return
	}

	// Display transformed data if no error occurred
	fmt.Println("Transformed data:", data)
}

// Function to handle different error types and perform rollback
func handleError(data Data, rollbackMap RollbackMap, err error) {
	switch err.(type) {
	case *ValidationError:
		ve := err.(*ValidationError)
		fmt.Printf("Validation error: %s\n", ve.Error())
		rollback(data, rollbackMap)
	case *TransformationError:
		te := err.(*TransformationError)
		fmt.Printf("Transformation error: %s\n", te.Error())
		rollback(data, rollbackMap)
	default:
		fmt.Printf("Unexpected error: %v\n", err)
		rollback(data, rollbackMap)
	}
}
