package main

import (
	"fmt"
	"log"
	"time"
)

// A simple callback function type
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

func harmfulCallback(msg string) {
	log.Printf("Potentially harmful message: %s\n", msg)
	// Simulate a harmful action, e.g., excessive resource usage
	for i := 0; i < 1000000; i++ {
		// This is a placeholder for a potentially harmful operation
		// like accessing sensitive data or executing system commands
	}
}

func suspiciousCallback(msg string) {
	log.Printf("Suspicious message: %s\n", msg)
	// Simulate a suspicious action, e.g., sending data to an external server
	fmt.Println("Suspicious data sent to an external server!")
}

func main() {
	// Initialize logger with a custom prefix
	log.SetPrefix("app: ")

	// Registering multiple callbacks
	processEvent("Event A", benignCallback)
	processEvent("Event B", harmfulCallback)
	processEvent("Event C", suspiciousCallback)

	// Adding a delay to demonstrate sequential execution
	time.Sleep(2 * time.Second)

	// Summary log
	log.Println("All events have been processed.")
}
