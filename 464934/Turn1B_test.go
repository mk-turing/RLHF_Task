package main

import (
	"testing"
)

func addPrefix(strings []string, prefix string) []string {
	var prefixedStrings []string
	for _, str := range strings {
		prefixedStrings = append(prefixedStrings, prefix+str)
	}
	return prefixedStrings
}

func TestAddPrefix(t *testing.T) {
	testCases := []struct {
		inputStrings   []string
		prefix         string
		expectedOutput []string
	}{
		{[]string{"a", "b", "c"}, "test", []string{"testa", "testb", "testc"}},
		{[]string{"apple", "banana", "cherry"}, "fruit", []string{"fruitapple", "fruitbanana", "fruitcherry"}},
		{[]string{}, "x", []string{}},
		{[]string{"x"}, "y", []string{"yx"}},
	}
	for _, testCase := range testCases {
		actualOutput := addPrefix(testCase.inputStrings, testCase.prefix)
		if !equal(actualOutput, testCase.expectedOutput) {
			t.Errorf("addPrefix(%v, %q) = %v, want %v", testCase.inputStrings, testCase.prefix, actualOutput, testCase.expectedOutput)
		}
	}
}

func equal(a, b []string) bool {
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
