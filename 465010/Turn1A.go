package main

import (
	"fmt"
)

// This function calculates the factorial of a non-negative integer.
//
// Usage: factorial(n) -> int
// - Computes the product of all positive integers less than or equal to n.
// - If n is negative, returns 0.
// - Example: `factorial(5)` computes `1 * 2 * 3 * 4 * 5 = 120`.
//
// Example Usage:
// result := factorial(5)
// fmt.Printf("The factorial of 5 is %d\n", result)
func factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 || n == 1 {
		return 1
	}
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return result
}

func main() {
	result := factorial(5)
	fmt.Printf("The factorial of 5 is %d\n", result)
}
