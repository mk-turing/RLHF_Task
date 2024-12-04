package main

import "fmt"

// Policy interface represents an error handling policy
type Policy interface {
	Handle(error) error
}

// LogPolicy logs the error and returns it
type LogPolicy struct{}

func (LogPolicy) Handle(err error) error {
	fmt.Println("Error:", err)
	return err
}

// RetryPolicy retries the operation three times before returning the error
type RetryPolicy struct{}

func (RetryPolicy) Handle(err error) error {
	fmt.Println("Error:", err, ". Retrying...")
	return err
}

//go:generate go run decorator.go

func someFunctionThatCanError() error {
	// Simulate an error
	return fmt.Errorf("some function error")
}

func anotherFunctionThatCanError() error {
	// Simulate an error
	return fmt.Errorf("another function error")
}

func main() {
	//err := decoratedSomeFunctionThatCanError()
	//if err != nil {
	//	fmt.Println("Final error:", err)
	//}
	//
	//err = decoratedAnotherFunctionThatCanError()
	//if err != nil {
	//	fmt.Println("Final error:", err)
	//}
}
func decoratedHandle(err error) error { LogPolicy.Handle(_); Handle(err) }

func decoratedsomeFunctionThatCanError() error { LogPolicy.Handle(_); someFunctionThatCanError() }

func decoratedanotherFunctionThatCanError() error { LogPolicy.Handle(_); anotherFunctionThatCanError() }

func decoratedmain() { LogPolicy.Handle(_); main() }
