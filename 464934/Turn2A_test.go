package main

import (
	"testing"
)

func TestRangeDataTypes(t *testing.T) {

	// Testing range over a slice of strings
	testCasesStrings := []struct {
		input    []string
		prefix   string
		expected []string
	}{
		{[]string{"a", "b", "c"}, "test", []string{"testa", "testb", "testc"}},
		{[]string{"apple", "banana", "cherry"}, "fruit", []string{"fruitapple", "fruitbanana", "fruitcherry"}},
		{[]string{}, "x", []string{}},
		{[]string{"x"}, "y", []string{"yx"}},
	}
	for _, testCase := range testCasesStrings {
		actualOutput := stringSliceWithPrefix(testCase.input, testCase.prefix)
		if !equalStrings(actualOutput, testCase.expected) {
			t.Errorf("stringSliceWithPrefix(%v, %q) = %v, want %v", testCase.input, testCase.prefix, actualOutput, testCase.expected)
		}
	}

	// Testing range over a map
	testCasesMap := []struct {
		input    map[string]int
		expected map[string]int
	}{
		{map[string]int{"a": 1, "b": 2, "c": 3}, map[string]int{"a": 1, "b": 2, "c": 3}},
		{map[string]int{"x": 10}, map[string]int{"x": 10}},
		{map[string]int{}, map[string]int{}},
	}
	for _, testCase := range testCasesMap {
		iteratedMap := mapWithKeysAndValues(testCase.input)
		if !equalMaps(iteratedMap, testCase.expected) {
			t.Errorf("mapWithKeysAndValues(%v) = %v, want %v", testCase.input, iteratedMap, testCase.expected)
		}
	}

	// Testing range over a slice of structs
	testCasesStructs := []struct {
		input    []Person
		modified []Person
	}{
		{
			[]Person{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}},
			[]Person{{Name: "Alice", Age: 26}, {Name: "Bob", Age: 31}},
		},
		{[]Person{{Name: "Charles", Age: 22}}, []Person{{Name: "Charles", Age: 23}}},
	}
	for _, testCase := range testCasesStructs {
		ageThisYear(testCase.input)
		if !equalStructs(testCase.input, testCase.modified) {
			t.Errorf("ageThisYear(%v) = %v, want %v", testCase.input, testCase.input, testCase.modified)
		}
	}
}

// Helper function for slicing strings with prefix
func stringSliceWithPrefix(strings []string, prefix string) []string {
	var prefixedStrings []string
	for _, str := range strings {
		prefixedStrings = append(prefixedStrings, prefix+str)
	}
	return prefixedStrings
}

// Helper function for iterating over a map
func mapWithKeysAndValues(inputMap map[string]int) map[string]int {
	outputMap := make(map[string]int)
	for key, value := range inputMap {
		outputMap[key] = value
	}
	return outputMap
}

// Helper function to modify ages in a slice of structs
func ageThisYear(people []Person) {
	for i := range people {
		people[i].Age++
	}
}

// Helper functions to compare elements
func equalStrings(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func equalMaps(a, b map[string]int) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for key, val := range a {
		if b[key] != val {
			return false
		}
	}

	return true
}

func equalStructs(a, b []Person) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Name != b[i].Name || a[i].Age != b[i].Age {
			return false
		}
	}

	return true
}

// Sample struct to use for testing slices of structs
type Person struct {
	Name string
	Age  int
}
