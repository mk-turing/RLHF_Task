package main

import "fmt"

// Callback definition
type Callback func(string)

// Function that accepts a callback
func processData(data string, callback Callback) {
	// Processing data...
	// Simulate some time-consuming task
	// ...

	// Invoke the callback
	callback(data)
}

// Unsafe callback implementation - Executes harmful code
func unsafeCallback(data string) {
	fmt.Println("Executing harmful code...")
	// Perform dangerous operations here, e.g., deleting files, writing to sensitive locations
	fmt.Println("Harmful code executed successfully.")
}

// Safe callback implementation - Validates data before executing operations
func safeCallback(data string) {
	// Validate data before proceeding
	if data == "secret" {
		fmt.Println("Received secret data. Performing sensitive operations...")
		// Perform sensitive operations
	} else {
		fmt.Println("Invalid data received. Ignoring.")
	}
}

func main() {
	dataToProcess := "unsafeData"
	// Passing the unsafe callback - Potential vulnerability
	processData(dataToProcess, unsafeCallback)

	// Passing the safe callback - Mitigates potential vulnerabilities
	processData(dataToProcess, safeCallback)
}
