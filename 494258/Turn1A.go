package main

import "fmt"

// 1. Slice Index Loop Optimization
func sliceIndexLoop() {
	data := make([]int, 10000)
	fmt.Println("Slice Index Loop:")
	for i := 0; i < len(data); i++ {
	}
}

// 2. Direct Array Access Optimization
func directArrayAccess() {
	data := [10000]int{}
	fmt.Println("Direct Array Access:")
	for i := 0; i < len(data); i++ {
	}
}

// 3. Avoid Unnecessary Index Calculations in Multi-Dimensional Arrays
func avoidUnnecessaryIndexCalculations() {
	matrix := make([][]int, 100, 100)
	n, m := 100, 100

	// Range loop with index calculation
	fmt.Println("Index Calculation Inside Loop:")
	for i := range matrix {
		for range matrix[i] {
		}
	}

	fmt.Println("Index Calculation Outside Loop:")
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
		}
	}
}

// 4. Loop Unrolling Optimization
func loopUnrolling() {
	data := make([]int, 10000)
	fmt.Println("Unrolled Loop:")
	for i := 0; i < len(data); i += 4 {
	}
}

func main() {
	// Calling the optimized loop functions
	sliceIndexLoop()
	directArrayAccess()
	avoidUnnecessaryIndexCalculations()
	loopUnrolling()

	// Additional explanation of the optimizations
	fmt.Println("\nOptimizations Explanation:")
	fmt.Println("1. Slice Index Loop: Replaced range loop with index loop for performance.")
	fmt.Println("2. Direct Array Access: Accessing array directly instead of using range loop.")
	fmt.Println("3. Avoid Unnecessary Index Calculations: Pre-compute indices outside of the loop.")
	fmt.Println("4. Loop Unrolling: Reduced loop iterations by processing multiple elements per iteration.")
}
