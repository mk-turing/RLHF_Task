package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func addPrefix(strings []string, prefix string) []string {
	var prefixedStrings []string
	for _, str := range strings {
		prefixedStrings = append(prefixedStrings, prefix+str)
	}
	return prefixedStrings
}

func TestAddPrefixPerformance(t *testing.T) {
	// Generate a large dataset
	const numStrings = 100000
	strings := make([]string, numStrings)
	for i := range strings {
		strings[i] = fmt.Sprintf("item%d", i)
	}

	startTime := time.Now()

	// Process the large dataset using range
	prefixedStrings := addPrefix(strings, "prefix_")

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)
	t.Logf("Execution time: %s", executionTime)

	// Validate the first few and last few items in the prefixed strings slice
	// for brevity and to ensure performance test doesn't take long
	for i := range []int{0, numStrings - 1} {
		if prefixedStrings[i] != "prefix_item"+fmt.Sprintf("%d", i) {
			t.Errorf("Unexpected prefixed string at index %d, got %q", i, prefixedStrings[i])
		}
	}
}

func TestAddPrefixToGeneric(t *testing.T) {
	testCases := []struct {
		input          interface{}
		prefix         string
		expectedOutput interface{}
	}{
		{[]string{"a", "b", "c"}, "test", []string{"testa", "testb", "testc"}},
		{[]int{1, 2, 3}, "num", []int{1, 2, 3}},
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
		{nil, "x", nil},
		{[]string{}, "y", []string{}},
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
		var prefixedInts []int
		for _, num := range v {
			strNum := prefix + fmt.Sprint(num)
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
