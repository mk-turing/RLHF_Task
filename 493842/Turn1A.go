package main

import (
	"fmt"
)

// A simple callback function type
type Callback func(string)

// A function that accepts a callback and uses it
func processEvent(message string, callback Callback) {
	if callback != nil {
		callback(message)
	}
}

func main() {
	// A benign callback function
	benignCallback := func(msg string) {
		fmt.Println("Benign message:", msg)
	}

	// Registering the benign callback
	processEvent("Hello, World!", benignCallback)

	// An example of a potentially malicious callback
	maliciousCallback := func(msg string) {
		fmt.Println("Malicious message:", msg)
		// Simulate a harmful action, e.g., execute system commands or access sensitive data
		// In a real-world scenario, this could be a call to a function that deletes files or exposes data
	}

	// Registering the malicious callback
	processEvent("Hello, Hacker!", maliciousCallback)
}
