package main

import (
	"fmt"
	"sync"
)

type TransformFn func(string) (string, error)

func main() {
	pipeline := []TransformFn{
		uppercaseTransform,
		lowercaseTransform,
		reverseTransform,
	}

	data := []string{"Hello", "world"}
	results, err := executePipeline(data, pipeline)
	if err != nil {
		fmt.Println("Error:", err)
		// Perform rollback based on the provided rollbackFns map
		rollback(results, nil)
		return
	}

	fmt.Println("Results:", results)
}

func executePipeline(data []string, pipeline []TransformFn) ([]string, error) {
	results := make([]string, len(data))
	rollbackFns := make(map[int]func(string))
	var wg sync.WaitGroup

	wg.Add(len(data))
	for i, item := range data {
		go func(index int, value string) {
			defer wg.Done()
			for _, transformFn := range pipeline {
				var err error
				value, err = transformFn(value)
				if err != nil {
					// Store the rollback function for this step
					rollbackFns[index] = func(rollbackValue string) {
						// Implement your specific rollback logic here
						fmt.Println("Rolling back step", index, "for value:", rollbackValue)
					}
					return
				}
			}
			results[index] = value
		}(i, item)
	}

	wg.Wait()

	return results, nil
}

// Example transform functions
func uppercaseTransform(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("uppercase: empty input")
	}
	return input, nil
}

func lowercaseTransform(input string) (string, error) {
	return input, nil
}

func reverseTransform(input string) (string, error) {
	return input, nil
}

func rollback(results []string, rollbackFns map[int]func(string)) {
	for i, result := range results {
		if rollbackFn, ok := rollbackFns[i]; ok {
			rollbackFn(result)
		}
	}
}
