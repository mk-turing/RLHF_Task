package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Custom Error Types
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Inner   error  `json:"-"`
}

// Error returns a string representation of the error
func (e *CustomError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("error code: %d, message: %s, inner error: %s", e.Code, e.Message, e.Inner.Error())
	}
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

// New creates a new CustomError
func NewCustomError(code int, message string, err error) error {
	return &CustomError{Code: code, Message: message, Inner: err}
}

// Error Handler Functions
func HandleError(err error) {
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

// RecoverError recovers from a panic and logs the error
func RecoverError() {
	if r := recover(); r != nil {
		log.Printf("recovered: %v", r)
	}
}

// Service Module
func ReadFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, NewCustomError(500, "failed to read file", err)
	}
	return data, nil
}

func ProcessFile(filePath string) error {
	data, err := ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error processing file: %w", err)
	}

	// Simulate processing logic
	if string(data) == "" {
		return NewCustomError(400, "file is empty", nil)
	}

	return nil
}

// Main Application Entry Point
func main() {
	defer RecoverError()

	filePath := filepath.Join("path", "to", "your", "file.txt")

	err := ProcessFile(filePath)
	HandleError(err)
}
