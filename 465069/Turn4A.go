package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 7 {
		fmt.Println("Usage: go run main.go <name1> <age1> <balance1> <name2> <age2> <balance2>")
		os.Exit(1)
	}

	// Extract user inputs
	name1 := os.Args[1]
	age1 := os.Args[2]
	balance1 := os.Args[3]
	name2 := os.Args[4]
	age2 := os.Args[5]
	balance2 := os.Args[6]

	// Define column widths
	nameWidth := 20
	ageWidth := 10
	balanceWidth := 15

	// Create table header
	fmt.Printf("| %-*s | %-*s | %-*s |\n", nameWidth, "Name", ageWidth, "Age", balanceWidth, "Balance")
	fmt.Printf("| %-*s | %-*s | %-*s |\n", nameWidth, "-", ageWidth, "-", balanceWidth, "-")

	// Create table rows
	fmt.Printf("| %-*s | %-*s | %-*s |\n", nameWidth, name1, ageWidth, age1, balanceWidth, balance1)
	fmt.Printf("| %-*s | %-*s | %-*s |\n", nameWidth, name2, ageWidth, age2, balanceWidth, balance2)
}
