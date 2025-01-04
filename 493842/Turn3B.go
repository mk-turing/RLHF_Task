package main

import (
	"fmt"
	"log"
	"time"
)

type Callback func(string)

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
	// Simulate harmful action
	for i := 0; i < 1000000; i++ {
	}
}

func suspiciousCallback(msg string) {
	log.Printf("Suspicious message: %s\n", msg)
	fmt.Println("Suspicious data sent to an external server!")
}

func main() {
	log.SetPrefix("app: ")

	// Registering predefined callbacks
	processEvent("Event A", benignCallback)
	processEvent("Event B", harmfulCallback)
	processEvent("Event C", suspiciousCallback)

	// Adding a delay
	time.Sleep(2 * time.Second)

	// User-defined callback input
	userInput := getUserInput("Enter a callback function (e.g., benignCallback): ")
	callback, err := getCallbackFromString(userInput)
	if err != nil {
		log.Fatalf("Error getting callback: %v", err)
	}

	// Process user-defined event
	processEvent("User Event", callback)

	// Adding a delay to demonstrate sequential execution
	time.Sleep(2 * time.Second)

	log.Println("All events have been processed.")
}

// Secure method to validate and restrict user-provided callbacks
func getCallbackFromString(userInput string) (Callback, error) {
	// Define a whitelist of allowed callback functions
	allowedCallbacks := map[string]Callback{
		"benignCallback": benignCallback,
	}

	// Convert user input to lowercase for case-insensitive comparison
	lowerInput := userInput

	// Check if the user input matches any of the allowed callback names
	for name, callback := range allowedCallbacks {
		if lowerInput == name {
			return callback, nil
		}
	}

	// If no match is found, return an error
	return nil, fmt.Errorf("invalid callback function: %s", userInput)
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
