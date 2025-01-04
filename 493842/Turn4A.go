package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// A callback function type that interacts with a file
type FileCallback func(filePath string)

func processFileEvent(filePath string, callback FileCallback) {
	if callback != nil {
		log.Printf("Executing callback for file: %s\n", filePath)
		callback(filePath)
	}
}

func benignFileCallback(filePath string) {
	log.Printf("Benign operation: Reading %s\n", filePath)
	// Simulate reading the file
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
	}
}

func maliciousFileCallback(filePath string) {
	log.Printf("Malicious operation: Writing to %s\n", filePath)
	// Simulate writing harmful data to the file
	_, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
	}
}

func protectedFileCallback(filePath string) {
	log.Printf("Protected operation: Reading %s\n", filePath)
	// Simulate reading the file
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
	}
}

func main() {
	log.SetPrefix("app: ")

	// Simulated file path
	filePath := "shared_resource.txt"

	// Registering predefined callbacks
	processFileEvent(filePath, benignFileCallback)
	processFileEvent(filePath, protectedFileCallback)

	// Delay for sequential execution
	time.Sleep(2 * time.Second)

	// Simulated user-defined callback input
	userInput := getUserInput("Enter a callback function (e.g., benignFileCallback): ")
	callback, err := getCallbackFromString(userInput)
	if err != nil {
		log.Fatalf("Error getting callback: %v", err)
	}

	// Process user-defined event
	processFileEvent(filePath, callback)

	// Delay for sequential execution
	time.Sleep(2 * time.Second)

	log.Println("All file events have been processed.")
}

// Secure method to validate and restrict user-provided callbacks
func getCallbackFromString(userInput string) (FileCallback, error) {
	// Define a whitelist of allowed callback functions
	allowedCallbacks := map[string]FileCallback{
		"benignFileCallback":    benignFileCallback,
		"protectedFileCallback": protectedFileCallback,
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
