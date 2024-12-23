package main

import (
	"fmt"
	"os"
)

func main() {
	// Command-line arguments
	if len(os.Args) != 3 {
		// Use fmt.Sprintf to format an error message
		errorMessage := fmt.Sprintf("Error: Missing required arguments. Please provide two arguments, e.g., %s <first_name> <last_name>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, errorMessage)
		os.Exit(1)
	}

	firstName := os.Args[1]
	lastName := os.Args[2]

	// Use fmt.Sprintf to create a formatted greeting
	greeting := fmt.Sprintf("Hello, %s %s! Welcome to our CLI app.", firstName, lastName)

	// Display the formatted greeting
	fmt.Println(greeting)
}
