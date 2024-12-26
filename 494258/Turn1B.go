package main

import "fmt"

// 1. Avoid Redundant Calculations
func avoidRedundantCalculations() {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	fmt.Println(sum) // Output: 499999500000
}

// 2. Use Appropriate Data Types
func useAppropriateDataTypes() {
	var count int = 0 // Using 'int' instead of 'int32'
	for count < 1000000 {
		count++
	}
	fmt.Println(count) // Output: 1000000
}

// 3. Use break and continue Statements Wisely
func useBreakAndContinueWisely() {
	// Using 'break' to stop the loop when 'i' is 50
	for i := 0; i < 100; i++ {
		if i == 50 {
			break
		}
		fmt.Println(i)
	}

	// Using 'continue' to skip the iteration when 'i' is 50
	for i := 0; i < 100; i++ {
		if i == 50 {
			continue
		}
		fmt.Println(i)
	}
}

// 4. Avoid Unnecessary Indexing
func avoidUnnecessaryIndexing() {
	data := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, value := range data { // For-range loop without unnecessary indexing
		sum += value
	}
	fmt.Println(sum) // Output: 15
}

// 5. Use Parallelism (incomplete)
func useParallelism() {
	// Goroutines and channels can be used here for parallel processing
	// This part is left incomplete for future implementation.
}

func main() {
	// Calling the optimized loop functions
	avoidRedundantCalculations()
	useAppropriateDataTypes()
	useBreakAndContinueWisely()
	avoidUnnecessaryIndexing()
	useParallelism() // This part is incomplete

	fmt.Println("\nOptimizations Explanation:")
	fmt.Println("1. Avoid Redundant Calculations: Calculated values are stored in variables to avoid repeated calculations.")
	fmt.Println("2. Use Appropriate Data Types: Using 'int' instead of 'int32' for performance improvement.")
	fmt.Println("3. Use break and continue Statements Wisely: Avoid excessive control flow changes.")
	fmt.Println("4. Avoid Unnecessary Indexing: For-range loops used to avoid unnecessary index-based access.")
	fmt.Println("5. Use Parallelism: Parallelization left for future implementation.")
}
