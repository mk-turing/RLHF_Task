package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// ErrorWithContext is a generic error that wraps another error with additional context.
type ErrorWithContext[T any] struct {
	error
	context interface{} // Changed to interface{} to hold any type, including string
}

func (e ErrorWithContext[T]) Error() string {
	// Convert context to string if it is of type string, otherwise use %+v for complex types
	if ctx, ok := e.context.(string); ok {
		return fmt.Sprintf("%v: %v", ctx, e.error)
	} else {
		return fmt.Sprintf("%+v: %v", e.context, e.error)
	}
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

// logResult logs a result, including errors.
func logResult[T any](result Result[T], level string) {
	if result.Error != nil {
		// Log error with context if the error is an ErrorWithContext
		if e, ok := result.Error.(ErrorWithContext[T]); ok {
			log.Printf("%s: Error: %s\n", level, e.Error())
		} else {
			log.Printf("%s: Error: %s\n", level, result.Error.Error())
		}
	} else {
		// Log the value
		log.Printf("%s: Result: %s\n", level, reflect.TypeOf(result.Value).String())
	}
}

// existsWithError checks if an element exists in a collection and returns an error if it doesn't.
func existsWithError[T comparable](col []T, target T) (Result[bool], error) {
	for _, element := range col {
		if element == target {
			return Result[bool]{Value: true, Error: nil}, nil
		}
	}
	return Result[bool]{}, fmt.Errorf("element %v not found", target)
}

// CheckList includes type constraints for `T` to be comparable and uses existsWithError.
func CheckList[T comparable](list []T, elements []T) (Result[[]T], error) {
	for _, element := range elements {
		_, err := existsWithError[T](list, element)
		if err != nil {
			// Pass the formatted string as part of the context (now it can be any type)
			wrappedError := ErrorWithContext[T]{
				error:   err,
				context: fmt.Sprintf("list (%v) does not contain element %v", list, element), // context is now a string
			}
			return Result[[]T]{}, wrappedError
		}
	}
	return Result[[]T]{Value: list, Error: nil}, nil
}

func main() {
	// Example usage
	inputList := []string{"apple", "banana", "cherry"}
	elementsToFind := []string{"banana", "date"}

	result, err := CheckList[string](inputList, elementsToFind)

	logResult(result, "INFO")
	if err != nil {
		wrappedError, ok := err.(ErrorWithContext[string])
		if ok {
			log.Println("Detailed Error:", wrappedError.Error())
		} else {
			log.Println("General Error:", err)
		}
	}
}
