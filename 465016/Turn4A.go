package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
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

	// Rollback all stages in reverse order with synchronization
	rollbackMutex := &sync.Mutex{}
	for i := len(stages) - 1; i >= 0; i-- {
		rollbackStage := stages[i]
		rollbackMutex.Lock()
		defer rollbackMutex.Unlock()
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

// Run transformations in parallel
func runTransformsInParallel(data Data, transformations []func(Data, *map[string]interface{}, string) error) ([]map[string]interface{}, error) {
	var stages []map[string]interface{}
	var err error

	// Create a channel for errors
	errorsCh := make(chan error)
	var wg sync.WaitGroup

	// Launch goroutines for each transformation
	for _, transform := range transformations {
		rollbackStage := make(map[string]interface{})
		wg.Add(1)
		go func(transform func(Data, *map[string]interface{}, string) error, rollbackStage *map[string]interface{}, errorsCh chan error, wg *sync.WaitGroup) {
			defer wg.Done()
			for key := range data {
				if err := transform(data, rollbackStage, key); err != nil {
					return
				}
			}
			stages = append(stages, *rollbackStage)
		}(transform, &rollbackStage, errorsCh, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return stages, err
}

func main() {
	data := Data{
		"a": 10,
		"b": "hello",
		"c": "world",
	}

	// Define transformations to run in parallel
	transformations := []func(Data, *map[string]interface{}, string) error{
		transformMultiply,
		transformStringToUpper,
	}

	// Run transformations in parallel
	stages, err := runTransformsInParallel(data, transformations)
	if err != nil {
		handleError(data, stages, err)
		return
	}

	// Display transformed data if no error occurred
	fmt.Println("Transformed data:", data)
}
