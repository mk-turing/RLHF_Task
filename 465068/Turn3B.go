package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type ErrorWithContext[T any] struct {
	error
	context T
}

func (e ErrorWithContext[T]) Error() string {
	return fmt.Sprintf("%+v: %v", e.context, e.error)
}

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

// logResult logs a result, including errors, with a general type constraint
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

// Define MyError outside the main function
type MyError struct {
	message string
}

func (me MyError) Error() string {
	return me.message
}

func main() {
	// Example usage
	input := "valid"
	result, _ := processData(input)

	logResult(result, "INFO")

	// Example with an error
	input = "error"
	result, _ = processData(input)

	logResult(result, "INFO")

	// Example with custom error type
	myError := MyError{message: "my custom error"}
	logResult(Result[error]{Error: myError}, "INFO")
}
