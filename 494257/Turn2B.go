package main

import (
	"fmt"
	"math"
)

const (
	maxAbsDifference = 1.0
)

func isValid(data, pattern []float64, index int) bool {
	for i := 0; i < len(pattern); i++ {
		if data[index+i] == math.NaN() || math.Abs(data[index+i]-pattern[i]) > maxAbsDifference {
			return false
		}
	}
	return true
}

func main() {
	data := []float64{1.0, 2.0, 3.0, math.NaN(), 5.0, 6.0, math.NaN(), 8.0, 9.0, 10.0, 9.0, 8.0, 7.0, 6.0, math.NaN(), 4.0, 3.0, 2.0, 1.0}
	pattern := []float64{2.0, 3.0, 4.0, 5.0}
	match := false

	for i := 0; i < len(data)-len(pattern)+1; i++ {
		match = true
		if isValid(data, pattern, i) {
			for j, v := range pattern {
				if data[i+j] != v {
					match = false
					break
				}
			}
			if match {
				fmt.Println("Pattern found at index", i)
			}
		}
	}
}
