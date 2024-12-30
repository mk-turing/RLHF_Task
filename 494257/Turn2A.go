package main

import (
	"fmt"
	"math"
)

func main() {
	// Sample scientific data with missing values and noise
	data := []float64{1.0, 2.0, 0.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 9.1, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0}
	pattern := []float64{2.0, 3.0, 4.0, 5.0}
	matchThreshold := 0.05 // Tolerance for noise

	// Using range loop for pattern matching
	for i := 0; i < len(data)-len(pattern)+1; i++ {
		match := true
		for j, v := range pattern {
			if data[i+j] == 0.0 {
				// Skip missing values
				continue
			}
			// Calculate absolute difference to account for noise
			diff := math.Abs(data[i+j] - v)
			if diff > matchThreshold {
				match = false
				break
			}
		}
		if match {
			fmt.Println("Pattern found at index", i)
		}
	}
}
