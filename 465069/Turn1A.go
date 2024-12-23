package main

import (
	"fmt"
)

func main() {
	// Define some sample data
	name := "Alice"
	age := 30
	email := "alice@example.com"

	// Use fmt.Sprintf to format a string with placeholders
	output := fmt.Sprintf("User Information:\nName: %s\nAge: %d\nEmail: %s\n", name, age, email)

	// Print the formatted output
	fmt.Println(output)
}
