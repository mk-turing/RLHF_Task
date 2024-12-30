package main

import (
	"fmt"
	"math"
)

const (
	maxAbsDifference = 1.0
	numWorkers       = 4 // Number of worker goroutines
)

func isValid(data, pattern []float64, index int) bool {
	for i := 0; i < len(pattern); i++ {
		if data[index+i] == math.NaN() || math.Abs(data[index+i]-pattern[i]) > maxAbsDifference {
			return false
		}
	}
	return true
}

func patternMatcher(data, pattern []float64, start, end int, result chan<- int) {
	for i := start; i <= end; i++ {
		if isValid(data, pattern, i) {
			result <- i
		}
	}
}

func main() {
	data := []float64{1.0, 2.0, 3.0, math.NaN(), 5.0, 6.0, math.NaN(), 8.0, 9.0, 10.0, 9.0, 8.0, 7.0, 6.0, math.NaN(), 4.0, 3.0, 2.0, 1.0}
	pattern := []float64{2.0, 3.0, 4.0, 5.0}
	result := make(chan int)

	// Divide the work into chunks
	chunkSize := (len(data) - len(pattern) + numWorkers - 1) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := min((i+1)*chunkSize-1, len(data)-len(pattern))
		go patternMatcher(data, pattern, start, end, result)
	}

	// Collect results
	for i := 0; i < numWorkers; i++ {
		if index, ok := <-result; ok {
			fmt.Println("Pattern found at index", index)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
