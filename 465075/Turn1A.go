package main

import (
	"fmt"
)

func main() {
	var name *string = nil
	var age int = 30
	var location *string = "New York"

	if name == nil {
		name = "Anonymous"
	}

	if location == nil {
		location = "Unknown"
	}

	fmt.Printf("Hello, my name is %s. I am %d years old and I live in %s.\n", name, age, location)

	// Using if let for safer handling
	if name_, ok := name.(*string); ok {
		fmt.Printf("Hello, my name is %s.\n", name_)
	} else {
		fmt.Printf("Hello, my name is Anonymous.\n")
	}

	if location_, ok := location.(*string); ok {
		fmt.Printf("I live in %s.\n", location_)
	} else {
		fmt.Printf("I live in Unknown.\n")
	}
}
