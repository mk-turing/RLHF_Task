package main

import (
	"fmt"
	"math"
)

// FunctionDetails is a struct that holds information about a function.
type FunctionDetails struct {
	Name      string
	Params    []string
	Returns   []string
}

// FormatFunctionComment uses fmt.Sprintf to format a structured comment for a function.
func FormatFunctionComment(details FunctionDetails) string {
	return fmt.Sprintf(`
	Function: %s
	Parameters:
	%s
	Returns:
	%s
	`, details.Name,
		fmt.Sprintf("%s\n",
			fmt.Sprintf("\t%s\n", p) for p in details.Params),
	fmt.Sprintf("%s\n",
		fmt.Sprintf("\t%s\n", r) for r in details.Returns),
)
}

// CalculateArea calculates the area of a circle.
func CalculateArea(radius float64) float64 {
	return math.Pi * radius * radius
}

// Main function to demonstrate the usage of FormatFunctionComment.
func main() {
	// Function details for CalculateArea.
	calculateAreaDetails := FunctionDetails{
		Name:      "CalculateArea",
		Params:    []string{"radius float64"},
		Returns:   []string{"float64"},
	}

	// Print the formatted comment.
	fmt.Println(FormatFunctionComment(calculateAreaDetails))

	// Call the CalculateArea function.
	result := CalculateArea(5.0)
	fmt.Println("Area:", result)
}