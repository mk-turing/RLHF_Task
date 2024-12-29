package main

import (
	"fmt"
	"testing"
)

// The function that takes a callback and performs an operation
func performOperation(callback func(int) error) error {
	num := 10
	err := callback(num)
	if err != nil {
		return err
	}
	return nil
}

// A mock callback function for testing purposes
func mockCallback(num int) error {
	if num == 10 {
		return fmt.Errorf("Error occurred with number %d", num)
	}
	return nil
}

func TestPerformOperation_WithMockCallback(t *testing.T) {
	// Test 1: Verify that the function returns an error when the mock callback returns an error
	err := performOperation(mockCallback)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	// Test 2: Verify that the function returns no error when the mock callback returns no error
	mockCallbackNoError := func(num int) error {
		return nil
	}
	err = performOperation(mockCallbackNoError)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestPerformOperation_WithNilCallback(t *testing.T) {
	// Test 3: Verify that the function returns an error when the callback is nil
	err := performOperation(nil)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}
