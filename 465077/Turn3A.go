package main

import (
	"fmt"
	"io"
)

type Person struct {
	Name string
	Age  int
}

type Address struct {
	Street string
	City   string
}

func (p Person) Format(f io.Writer) {
	fmt.Fprintf(f, "Person{Name: %q, Age: %d}", p.Name, p.Age)
}

func (a Address) Format(f io.Writer) {
	fmt.Fprintf(f, "Address{Street: %q, City: %q}", a.Street, a.City)
}

func main() {
	person := Person{Name: "Alice", Age: 25}
	address := Address{Street: "Main Street", City: "New York"}
	nestedData := map[string]interface{}{
		"person":  person,
		"address": address,
	}

	// Using %v with the Formatter interface
	fmt.Printf("Nested Data: %v\n", nestedData) // Output: Nested Data: map[person:Person{Name: "Alice", Age: 25} address:Address{Street: "Main Street", City: "New York"}]
}
