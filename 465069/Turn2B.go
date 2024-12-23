package main

import (
	"fmt"
	"os"
)

func main() {
	// Command-line arguments
	args := os.Args[1:]

	if len(args) < 1 {
		// Format and print the error message
		errMsg := fmt.Sprintf("Error: Missing required argument <file_name>. Usage: go run main.go <file_name>")
		fmt.Println(errMsg)
		os.Exit(1)
	}

	fileName := args[0]
	fmt.Println("Processing file:", fileName)

	// Rest of the application logic
}
