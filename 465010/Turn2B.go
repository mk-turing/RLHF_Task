package main

import "fmt"

func calculateTotal(num1 int, num2 int) int {
	total := num1 + num2

	// Document the input values and the calculated total using fmt.Sprintf
	fmt.Println(fmt.Sprintf(`
		Calculation Details:
			Input 1: %d
			Input 2: %d
			Total:   %d
		`, num1, num2, total))

	return total
}

func main() {
	result := calculateTotal(10, 20)
	fmt.Println("Result:", result)
}
