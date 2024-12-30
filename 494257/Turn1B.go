package main

import (
	"fmt"
)

func main() {

	//Sample scientific data
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0}
	pattern := []float64{2.0, 3.0, 4.0, 5.0}
	match := false

	//Using range loop for pattern matching
	for i := 0; i < len(data)-len(pattern)+1; i++ {
		match = true
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

	//Time complexity analysis
	//Time complexity of this algorithm is O(n*m) where n is the length of the data array and m is the length of the pattern array.
	//This is because we are using nested loops to compare each element of the pattern with corresponding elements of the data array.
}
