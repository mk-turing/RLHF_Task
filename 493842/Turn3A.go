package main

import (
	"log"
	"time"
)

type Callback func(string)

// A function that accepts a callback and uses it
func processEvent(message string, callback Callback) {
	if callback != nil {
		log.Printf("Executing callback for event: %s\n", message)
		callback(message)
	}
}

func benignCallback(msg string) {
	log.Printf("Benign message: %s\n", msg)
}

// Compiles and executes user-provided code as a callback function
func compileCallback(callbackStr string) (Callback, error) {
	// Placeholder code to simulate compiling and executing a callback
	log.Printf("Compiling user-provided callback: %s\n", callbackStr)
	return nil, nil // In a real implementation, this would use go/parser and go/ast to compile code
}

func main() {
	log.SetPrefix("app: ")

	// Registering a benign callback
	processEvent("Event A", benignCallback)

	// Simulate user-provided callbacks
	unsafeCallbackStr := "func(msg string) { log.Printf(\"Unsafe message: %s\", msg); for i := 0; i < 1000000; i++ {} }"
	log.Printf("Received user-provided callback: %s\n", unsafeCallbackStr)

	// Compile and execute the user-provided callback
	unsafeCallback, err := compileCallback(unsafeCallbackStr)
	if err != nil {
		log.Fatalf("Error compiling callback: %v\n", err)
	}
	processEvent("Event B", unsafeCallback)

	// Adding a delay to demonstrate sequential execution
	time.Sleep(2 * time.Second)

	// Summary log
	log.Println("All events have been processed.")
}
