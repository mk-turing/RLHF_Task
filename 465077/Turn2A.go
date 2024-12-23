package main

import (
	"fmt"
	"strings"
)

type Person struct {
	Name string
	Age  int
}

type Address struct {
	Street string
	City   string
}

func formatNestedMap(m map[interface{}]interface{}, indent string) string {
	var result strings.Builder
	result.WriteString("{\n")
	for k, v := range m {
		result.WriteString(indent)
		result.WriteString(fmt.Sprintf("%#v: ", k))
		if s, ok := v.(string); ok {
			result.WriteString(fmt.Sprintf("%q", s))
		} else if mapValue, ok := v.(map[interface{}]interface{}); ok {
			result.WriteString(formatNestedMap(mapValue, indent+"  "))
		} else if structValue, ok := v.(interface{ fmt.Stringer }); ok {
			result.WriteString(structValue.String())
		} else {
			result.WriteString(fmt.Sprintf("%#v", v))
		}
		result.WriteString(",\n")
	}
	//if result.Len() > 2 {
	//	result = result[:result.Len()-2]
	//}
	result.WriteString("\n")
	result.WriteString(indent)
	result.WriteString("}")
	return result.String()
}

func main() {
	nestedMap := map[interface{}]interface{}{
		"person":  &Person{Name: "John Doe", Age: 30},
		"address": &Address{Street: "Main Street", City: "New York"},
		"skills": map[interface{}]interface{}{
			"programming": "Go",
			"languages":   []string{"English", "Spanish"},
		},
	}

	fmt.Println(formatNestedMap(nestedMap, ""))
}
