package main

import (
	"fmt"
	"strconv"
	"strings"
)

type User struct {
	Name    string
	Age     int
	Balance float64
}

func main() {
	// Take user input for multiple users
	users := []User{}
	for {
		var name, ageStr, balanceStr string
		fmt.Print("Enter Name: ")
		fmt.Scanln(&name)
		if name == "" {
			break
		}
		fmt.Print("Enter Age: ")
		fmt.Scanln(&ageStr)
		fmt.Print("Enter Balance: ")
		fmt.Scanln(&balanceStr)

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			fmt.Println("Invalid age format. Please enter a number.")
			continue
		}
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err != nil {
			fmt.Println("Invalid balance format. Please enter a number.")
			continue
		}
		users = append(users, User{name, age, balance})
	}
	// Calculate the maximum length for each column to determine alignment
	maxNameLength := 0
	maxAgeLength := 0
	maxBalanceLength := 0
	for _, user := range users {
		maxNameLength = max(maxNameLength, len(user.Name))
		maxAgeLength = max(maxAgeLength, len(strconv.Itoa(user.Age)))
		maxBalanceLength = max(maxBalanceLength, len(strconv.FormatFloat(user.Balance, 'f', 2, 64)))
	}

	// Format and display the table
	fmt.Println("\nUser Information:")
	fmt.Printf("| %-*s | %-*s | %-*s |\n", maxNameLength, "Name", maxAgeLength, "Age", maxBalanceLength, "Balance")
	fmt.Println(strings.Repeat("-", maxNameLength+maxAgeLength+maxBalanceLength+7))
	for _, user := range users {
		fmt.Printf("| %-*s | %-*d | %-*.*f |\n", maxNameLength, user.Name, maxAgeLength, user.Age, maxBalanceLength, 2, user.Balance)
	}
	fmt.Println(strings.Repeat("-", maxNameLength+maxAgeLength+maxBalanceLength+7))
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
