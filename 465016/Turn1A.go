package main

import (
	"errors"
	"fmt"
)

type Data map[string]interface{}
type RollbackMap map[string]interface{}

// Example transformation step
func transformAdd(data Data, rollbackMap *RollbackMap, key string, value interface{}) error {
	if data[key] == nil {
		data[key] = value
		(*rollbackMap)[key] = nil
	} else {
		return fmt.Errorf("Key %s already exists", key)
	}
	return nil
}

// Example transformation step
func transformMultiply(data Data, rollbackMap *RollbackMap, key string) error {
	if value, ok := data[key].(int); ok {
		data[key] = value * 2
		(*rollbackMap)[key] = value
	} else {
		return errors.New("Invalid type for multiplication")
	}
	return nil
}

// Rollback the changes recorded in the rollback map
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
		fmt.Printf("Error adding new key: %v\n", err)
		rollback(data, rollbackMap)
		return
	}

	err = transformMultiply(data, &rollbackMap, "c")
	if err != nil {
		fmt.Printf("Error multiplying value: %v\n", err)
		rollback(data, rollbackMap)
		return
	}

	// Display transformed data if no error occurred
	fmt.Println("Transformed data:", data)
}
