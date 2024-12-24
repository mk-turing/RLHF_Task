package main

import (
	"errors"
	"fmt"
)

type Data map[string]interface{}
type RollbackMap map[string]interface{}
type RollbackOperation struct {
	Key   string
	Value interface{}
	Action string // e.g., "add", "multiply", "delete"
}

// Define custom error types
type ValidationError struct {
	Key string
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for key '%s': %v", e.Key, e.Err)
}

type TransformationError struct {
	Key string
	Err error
}

func (e *TransformationError) Error() string {
	return fmt.Sprintf("Transformation error for key '%s': %v", e.Key, e.Err)
}

// Example transformation step: Add a key-value pair
func transformAdd(data Data, rollbackStack *[]RollbackOperation, key string, value interface{}) error {
	if data[key] != nil {
		return &ValidationError{Key: key, Err: fmt.Errorf("Key %s already exists", key)}
	}
	data[key] = value
	rollbackOp := RollbackOperation{Key: key, Value: nil, Action: "delete"}
	*rollbackStack = append(*rollbackStack, rollbackOp)
	return nil
}

// Example transformation step: Multiply the value of an existing key
func transformMultiply(data Data, rollbackStack *[]RollbackOperation, key string) error {
	value, ok := data[key].(int)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for multiplication")}
	}
	oldValue := value
	data[key] = value * 2
	rollbackOp := RollbackOperation{Key: key, Value: oldValue, Action: "multiply"}
	*rollbackStack = append(*rollbackStack, rollbackOp)
	return nil
}

// Example transformation step: Convert string value to uppercase
func transformUppercase(data Data, rollbackStack *[]RollbackOperation, key string) error {
	value, ok := data[key].(string)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for uppercase conversion")}
	}
	oldValue := value
	data[key] = fmt.Sprintf("%s", value)
	rollbackOp := RollbackOperation{Key: key, Value: oldValue, Action: "uppercase"}
	*rollbackStack = append(*rollbackStack, rollbackOp)
	return nil
}

// Perform rollback operations based on the stack
func rollback(data Data, rollbackStack []RollbackOperation) {
	for len(rollbackStack) > 0 {
		op := rollbackStack[len(rollbackStack)-1]
		rollbackStack = rollbackStack[:len(rollbackStack)-1]
		switch op.Action {
		case "add":
			delete(data, op.Key)
		case "multiply":
			data[op.Key] = op.Value
		case "uppercase":
			data[op.Key] = op.Value
		case "delete":
			data[op.Key] = op.Value
		}
	}
}

func main() {
	data := Data{
		"a": 10,
		"b": "hello",
	}
	var rollbackStack []RollbackOperation

	// Transaction starts
	err := transformAdd(data, &rollbackStack, "c", 20)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		rollback(data, rollbackStack)
		return
	}
	err = transformMultiply(data, &rollbackStack, "a")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		rollback(data, rollbackStack)