package main

import (
	"errors"
	"fmt"
)

// ErrorWithContext is a generic error that wraps another error with additional context.
type ErrorWithContext[T any] struct {
	error
	context T
}

func (e ErrorWithContext[T]) Error() string {
	return fmt.Sprintf("%+v: %v", e.context, e.error)
}

// Result is a generic result type that contains a value and an error.
type Result[T any] struct {
	Value T
	Error error
}

func createValueWithError[T any](value T, err error) Result[T] {
	return Result[T]{
		Value: value,
		Error: err,
	}
}

func processData[T any](data T) (Result[T], error) {
	// Simulating an error condition
	if str, ok := any(data).(string); ok && str == "error" {
		return Result[T]{}, errors.New("simulated error")
	}

	// Simulating successful processing
	return Result[T]{
		Value: data,
		Error: nil,
	}, nil
}

// logger is a generic logger function that handles different types of errors and values.
func logger[T any](result Result[T]) {
	if result.Error != nil {
		// Use the ErrorWithContext if available
		if errWithContext, ok := result.Error.(ErrorWithContext[any]); ok {
			fmt.Println("Error:", errWithContext.Error())
		} else {
			fmt.Println("Error:", result.Error.Error())
		}
	} else {
		fmt.Println("Result:", result.Value)
	}
}

func main() {
	// Example usage
	input := "valid"
	result, _ := processData(input)
	logger(result) // Output: Result: valid

	// Example with an error
	input = "error"
	result, _ = processData(input)
	logger(result) // Output: Error: simulated error

	// Example with a custom error with context
	customErr := errors.New("custom error")
	wrappedError := ErrorWithContext[string]{
		error:   customErr,
		context: "Some custom context",
	}
	logger(Result[string]{Error: wrappedError}) // Output: Error: Some custom context: custom error
}
