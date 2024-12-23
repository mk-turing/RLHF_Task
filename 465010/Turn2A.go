package main

import (
	"fmt"
)

// Person represents a human being.
type Person struct {
	Name    string
	Age     int
	Address string
}

func main() {
	// person is an instance of the Person struct.
	person := Person{
		Name:    "Alice",
		Age:     30,
		Address: "New York City",
	}

	// Create a structured comment using fmt.Sprintf to document the Person instance.
	fmt.Printf(`
		// Person:
		// 	Name: %s
		// 	Age: %d
		// 	Address: %s
		`, person.Name, person.Age, person.Address)
}
