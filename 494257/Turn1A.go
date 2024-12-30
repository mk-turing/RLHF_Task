package main

import (
	"fmt"
	"reflect"
)

func main() {
	// Sample dataset of scientific measurements
	dataset := []float64{1.2, 3.4, 2.3, 5.6, 3.4, 1.2, 3.4, 7.8, 3.4, 2.3, 1.2}

	// Define the pattern you want to find
	pattern := []float64{3.4, 2.3}

	// The length of the pattern
	patternLength := len(pattern)

	// Variable to keep track of the pattern count
	patternCount := 0

	// Use a range loop to iterate over the dataset
	for i := range dataset {
		// Slice of the dataset to compare with the pattern
		subdata := dataset[i : i+patternLength]

		// Compare the slices
		if reflect.DeepEqual(subdata, pattern) {
			// If they match, increment the pattern count
			patternCount++
		}
	}

	// Output the result
	fmt.Printf("Pattern found %d times in the dataset.\n", patternCount)
}
