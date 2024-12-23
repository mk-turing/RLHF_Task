package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

type Address struct {
	Street string
	City   string
}

func main() {
	person := Person{Name: "John Doe", Age: 30}
	address := Address{Street: "Main Street", City: "New York"}
	nestedMap := map[string]interface{}{
		"person":  person,
		"address": address,
	}

	// Printing struct values and map values using %+v
	fmt.Println("Person:", person)        // Output: Person: {John Doe 30}
	fmt.Println("Address:", address)      // Output: Address: {Main Street New York}
	fmt.Println("Nested Map:", nestedMap) // Output: Nested Map: map[address:{Main Street New York} person:{John Doe 30}]

	// Printing struct values and map values using %#v
	fmt.Println("Person:", person)        // Output: Person: main.Person{Name:"John Doe", Age:30}
	fmt.Println("Address:", address)      // Output: Address: main.Address{Street:"Main Street", City:"New York"}
	fmt.Println("Nested Map:", nestedMap) // Output: Nested Map: map[string]interface {}{"address":main.Address{Street:"Main Street", City:"New York"}, "person":main.Person{Name:"John Doe", Age:30}}

	// Accessing nested values using nested %+v
	fmt.Println("Nested Person:", nestedMap["person"])   // Output: Nested Person: {John Doe 30}
	fmt.Println("Nested Address:", nestedMap["address"]) // Output: Nested Address: {Main Street New York}

	// Printing strings using %q
	str := "Hello, World!"
	fmt.Println("String:", str)                           // Output: String: Hello, World!
	fmt.Println("Quoted String:", fmt.Sprintf("%q", str)) // Output: Quoted String: "Hello, World!"

	// Printing nil values
	var nilPerson *Person
	fmt.Println("Nil Person:", nilPerson)                                // Output: Nil Person: <nil>
	fmt.Println("Nil Person (with %+v):", fmt.Sprintf("%+v", nilPerson)) // Output: Nil Person (with %+v): <nil>
}
