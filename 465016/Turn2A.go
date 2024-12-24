package main

import (
	"fmt"
)

type Data map[string]interface{}
type RollbackMap map[string]interface{}

type AdditionError struct{ key string }
type MultiplicationError struct{ key string }

func (e *AdditionError) Error() string {
	return fmt.Sprintf("Key %s already exists", e.key)
}

func (e *MultiplicationError) Error() string {
	return fmt.Sprintf("Invalid type for multiplication for key %s", e.key)
}

// Example transformation step
func transformAdd(data Data, rollbackMap *RollbackMap, key string, value interface{}) error {
	if data[key] == nil {
		data[key] = value
		(*rollbackMap)[key] = nil
	} else {
		return &AdditionError{key: key}
	}
	return nil
}

// Example transformation step
func transformMultiply(data Data, rollbackMap *RollbackMap, key string) error {
	if value, ok := data[key].(int); ok {
		data[key] = value * 2
		(*rollbackMap)[key] = value
	} else {
		return &MultiplicationError{key: key}
	}
	return nil
}

// Rollback the changes recorded in the rollback map
func rollback(data Data, rollbackMap RollbackMap) {
	for key, oldValue := range rollbackMap {
		if value, ok := oldValue.(int); ok {
			data[key] = value
		} else {
			delete(data, key)
		}
	}
	rollbackMap = nil
}

// Enhanced rollback function that handles different error types
func rollbackBasedOnError(data Data, rollbackMap RollbackMap, err error) {
	switch err := err.(type) {
	case *AdditionError:
		delete(data, err.key)
	case *MultiplicationError:
		if value, ok := data[err.key].(int); ok {
			data[err.key] = value / 2
		}
	default:
		fmt.Printf("Unhandled error type: %T\n", err)
	}

	// Always clear the rollback map after rollback
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
		fmt.Printf("Error adding new key: %v\n", err)
		rollbackBasedOnError(data, rollbackMap, err)
		return
	}

	err = transformMultiply(data, &rollbackMap, "c")
	if err != nil {
		fmt.Printf("Error multiplying value: %v\n", err)
		rollbackBasedOnError(data, rollbackMap, err)
		return
	}

	err = transformMultiply(data, &rollbackMap, "b")
	if err != nil {
		fmt.Printf("Error multiplying value: %v\n", err)
		rollbackBasedOnError(data, rollbackMap, err)
		return
	}

	// Display transformed data if no error occurred
	fmt.Println("Transformed data:", data)
}
