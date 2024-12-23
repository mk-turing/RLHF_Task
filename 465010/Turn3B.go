package main

import (
	"fmt"
)

// Add function adds two integers.
// The result of the addition is returned.
func Add(a, b int) int {
	// Document the function parameters and return value using fmt.Sprintf
	fmt.Println(fmt.Sprintf(`
		Function: mathutil.Add
		Parameters:
			a: %d
			b: %d
		Return Value: %d
		`, a, b, a+b))

	return a + b
}

func main() {
	sum := Add(15, 25)
	fmt.Println("Sum:", sum)
}
