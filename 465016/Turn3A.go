package main

import (
	"errors"
	"fmt"
	"strings"
)

type Data map[string]interface{}

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
func transformAdd(data Data, rollbackStage *map[string]interface{}, key string, value interface{}) error {
	if data[key] != nil {
		return &ValidationError{Key: key, Err: fmt.Errorf("Key %s already exists", key)}
	}
	data[key] = value
	(*rollbackStage)[key] = nil
	return nil
}

// Example transformation step
func transformMultiply(data Data, rollbackStage *map[string]interface{}, key string) error {
	value, ok := data[key].(int)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for multiplication")}
	}
	data[key] = value * 2
	(*rollbackStage)[key] = value
	return nil
}

// Example transformation step
func transformStringToUpper(data Data, rollbackStage *map[string]interface{}, key string) error {
	value, ok := data[key].(string)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for string conversion")}
	}
	data[key] = strings.ToUpper(value)
	(*rollbackStage)[key] = value
	return nil
}

// Function to handle different error types and perform rollback
func handleError(data Data, stages []map[string]interface{}, err error) {
	switch err.(type) {
	case *ValidationError:
		ve := err.(*ValidationError)
		fmt.Printf("Validation error: %s\n", ve.Error())
	case *TransformationError:
		te := err.(*TransformationError)
		fmt.Printf("Transformation error: %s\n", te.Error())
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}

	// Rollback all stages in reverse order
	for i := len(stages) - 1; i >= 0; i-- {
		rollbackStage := stages[i]
		for key, oldValue := range rollbackStage {
			switch v := oldValue.(type) {
			case int:
				data[key] = v
			case string:
				data[key] = v
			default:
				delete(data, key)
			}
		}
	}
}

func main() {
	data := Data{
		"a": 10,
		"b": "hello",
	}
	stages := make([]map[string]interface{}, 0)

	// Transform the data in stages
	rollbackStage := make(map[string]interface{})
	err := transformAdd(data, &rollbackStage, "c", 20)
	if err != nil {
		handleError(data, stages, err)
		return
	}
	stages = append(stages, rollbackStage)

	rollbackStage = make(map[string]interface{})
	err = transformMultiply(data, &rollbackStage, "c")
	if err != nil {
		handleError(data, stages, err)
		return
	}
	stages = append(stages, rollbackStage)

	rollbackStage = make(map[string]interface{})
	err = transformMultiply(data, &rollbackStage, "a")
	if err != nil {
		handleError(data, stages, err)
		return
	}
	stages = append(stages, rollbackStage)

	rollbackStage = make(map[string]interface{})
	err = transformStringToUpper(data, &rollbackStage, "b")
	if err != nil {
		handleError(data, stages, err)
		return
	}
	stages = append(stages, rollbackStage)

	// Display transformed data if no error occurred
	fmt.Println("Transformed data:", data)
}
