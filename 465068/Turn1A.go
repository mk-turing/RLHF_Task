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

func main() {
	// Example usage
	input := "valid"
	result, err := processData(input)

	if err != nil {
		// Use the ErrorWithContext to wrap the error with context
		wrappedError := ErrorWithContext[string]{
			error:   err,
			context: input,
		}
		fmt.Println("Error:", wrappedError.Error())
	} else {
		fmt.Println("Result:", result.Value)
	}

	// Example with an error
	input = "error"
	result, err = processData(input)

	if err != nil {
		// Use the ErrorWithContext to wrap the error with context
		wrappedError := ErrorWithContext[string]{
			error:   err,
			context: input,
		}
		fmt.Println("Error:", wrappedError.Error())
	} else {
		fmt.Println("Result:", result.Value)
	}
}
