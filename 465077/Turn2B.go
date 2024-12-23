package main

import (
	"os"
	"text/template"
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
	nestedData := map[string]interface{}{
		"person":  person,
		"address": address,
	}

	tmpl, err := template.New("nestedData").Parse("Name: {{.person.Name}} Age: {{.person.Age}} Address: Street: {{.address.Street}} City: {{.address.City}}")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, nestedData)
	if err != nil {
		panic(err)
	}
}
