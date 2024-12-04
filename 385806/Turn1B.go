//package main
//
//import (
//	"fmt"
//	"os"
//)
//
//// Define the PolicyError interface
//type PolicyError interface {
//	error
//	GetPolicy() string
//}
//
//// PanicsAnnotatedError implement PolicyError for errors that should panic
//type PanicsAnnotatedError struct {
//	err error
//}
//
//func (e *PanicsAnnotatedError) Error() string {
//	return e.err.Error()
//}
//
//func (e *PanicsAnnotatedError) GetPolicy() string {
//	return "panic"
//}
//
//// LogsAnnotatedError implement PolicyError for errors that should be logged
//type LogsAnnotatedError struct {
//	err error
//}
//
//func (e *LogsAnnotatedError) Error() string {
//	return e.err.Error()
//}
//
//func (e *LogsAnnotatedError) GetPolicy() string {
//	return "log"
//}
//
//// ReturnsAnnotatedError implement PolicyError for errors that should be returned
//type ReturnsAnnotatedError struct {
//	err error
//}
//
//func (e *ReturnsAnnotatedError) Error() string {
//	return e.err.Error()
//}
//
//func (e *ReturnsAnnotatedError) GetPolicy() string {
//	return "return"
//}
//
//// Function that "panics" on error
//func divideByZero(a, b int) (err error) {
//	defer func() {
//		if r := recover(); r != nil {
//			// If panic occurred, return a custom error with the panic info
//			if err1, ok := r.(error); ok {
//				err = &PanicsAnnotatedError{err: err1}
//			}
//		}
//	}()
//
//	result := a / b
//	fmt.Println(result)
//	return nil
//}
//
//// Function that "logs" on error
//func writeToFile(file string) error {
//	if err := os.WriteFile(file, nil, 0644); err != nil {
//		return &LogsAnnotatedError{err: err}
//	}
//	return nil
//}
//
//// Function that "returns" the error
//func readFromFile(file string) ([]byte, error) {
//	_, err := os.ReadFile(file)
//	if err != nil {
//		return nil, &ReturnsAnnotatedError{err: err}
//	}
//	return nil, nil
//}
//
//// PolicyHandler executes functions with the applied error handling policy
//func PolicyHandler(f func() error) {
//	err := f()
//	if err == nil {
//		return
//	}
//
//	switch err := err.(type) {
//	case *PanicsAnnotatedError:
//		panic(err.err)
//	case *LogsAnnotatedError:
//		fmt.Println("Error:", err.err)
//	case *ReturnsAnnotatedError:
//		fmt.Println("Function returned error:", err.err)
//	default:
//		fmt.Println("Unknown error type:", err)
//	}
//}
//
//func main() {
//	PolicyHandler(func() error {
//		err := divideByZero(1, 0)
//		return err // This will panic
//	}) // Will panic
//
//	if err := writeToFile("non_existent_file.txt"); err != nil {
//		fmt.Println("Writing failed:", err)
//	}
//
//	_, err := readFromFile("non_existent_file.txt")
//	if err != nil {
//		fmt.Println("Reading failed:", err)
//	}
//}
