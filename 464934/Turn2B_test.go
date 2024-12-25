package main

import (
	"fmt"
	"reflect"
	"strconv"
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
		input          interface{}
		prefix         string
		expectedOutput interface{}
	}{
		{[]string{"a", "b", "c"}, "test", []string{"testa", "testb", "testc"}},
		{[]int{1, 2, 3}, "num", []int{1, 2, 3}}, // Test with slice of int
		{map[string]string{"key1": "val1", "key2": "val2"}, "map", map[string]string{"mapkey1": "val1", "mapkey2": "val2"}},
		{[]struct {
			Name string
			Age  int
		}{
			{"Alice", 25},
			{"Bob", 30},
		}, "user", []struct {
			Name string
			Age  int
		}{
			{"userAlice", 25},
			{"userBob", 30},
		}},
		{nil, "x", nil},               // Test with nil input
		{[]string{}, "y", []string{}}, // Test with empty slice
	}

	for _, testCase := range testCases {
		actualOutput := addPrefixToGeneric(testCase.input, testCase.prefix)
		if !reflect.DeepEqual(actualOutput, testCase.expectedOutput) {
			t.Errorf("addPrefixToGeneric(%v, %q) = %v, want %v", testCase.input, testCase.prefix, actualOutput, testCase.expectedOutput)
		}
	}
}

func addPrefixToGeneric(input interface{}, prefix string) interface{} {
	switch v := input.(type) {
	case []string:
		return addPrefix(v, prefix)
	case []int:
		// Add prefix to ints by converting them to strings, adding prefix, and then converting back to ints (if required)
		var prefixedInts []int
		for _, num := range v {
			strNum := prefix + fmt.Sprint(num)
			// Handle conversion error (e.g., if prefix is non-numeric)
			if prefixedInt, err := strconv.Atoi(strNum); err == nil {
				prefixedInts = append(prefixedInts, prefixedInt)
			} else {
				return fmt.Errorf("cannot convert %q to int", strNum)
			}
		}
		return prefixedInts
	case map[string]string:
		prefixedMap := make(map[string]string)
		for key, val := range v {
			prefixedMap[prefix+key] = val
		}
		return prefixedMap
	case []struct {
		Name string
		Age  int
	}:
		prefixedStructs := make([]struct {
			Name string
			Age  int
		}, len(v))
		for i, structVal := range v {
			prefixedStructs[i].Name = prefix + structVal.Name
			prefixedStructs[i].Age = structVal.Age
		}
		return prefixedStructs
	default:
		return fmt.Errorf("unsupported input type: %T", input)
	}
}
