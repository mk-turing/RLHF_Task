package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

const (
	maxAbsDifference = 1.0
	windowSize       = 4
)

func isValid(data, pattern []float64, index int) bool {
	for i := 0; i < len(pattern); i++ {
		if data[index+i] == math.NaN() || math.Abs(data[index+i]-pattern[i]) > maxAbsDifference {
			return false
		}
	}
	return true
}

func findPattern(data []float64, pattern []float64, results chan<- int) {
	for i := 0; i < len(data)-windowSize+1; i++ {
		if isValid(data, pattern, i) {
			match := true
			for j, v := range pattern {
				if data[i+j] != v {
					match = false
					break
				}
			}
			if match {
				results <- i
			}
		}
	}
}

func main() {
	// Sample scientific data with multiple patterns to find
	data := []float64{1.0, 2.0, 3.0, math.NaN(), 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 9.0, 8.0, 7.0, 6.0, math.NaN(), 4.0, 3.0, 2.0, 1.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	pattern1 := []float64{2.0, 3.0, 4.0, 5.0}
	pattern2 := []float64{1.0, 2.0, 3.0, 4.0}

	results := make(chan int)
	wg := &sync.WaitGroup{}

	numCPU := runtime.NumCPU()
	fmt.Println("Using", numCPU, "CPUs")

	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			findPattern(data, pattern1, results)
			findPattern(data, pattern2, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for index := range results {
		fmt.Println("Pattern found at index", index)
	}
}
