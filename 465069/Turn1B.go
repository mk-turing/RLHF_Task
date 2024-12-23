package main

import (
	"fmt"
	"os"
)

func main() {
	// Command-line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <first_name> <last_name>")
		os.Exit(1)
	}

	firstName := os.Args[1]
	lastName := os.Args[2]

	// Use fmt.Sprintf to create a formatted string
	greeting := fmt.Sprintf("Hello, %s %s! Welcome to our CLI app.", firstName, lastName)

	// Display the formatted string
	fmt.Println(greeting)
}
