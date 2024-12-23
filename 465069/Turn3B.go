package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <first_argument> <second_argument>")
		os.Exit(1)
	}

	arg1 := os.Args[1]
	arg2 := os.Args[2]

	// Determine the types of the arguments
	isStringArg1 := true
	isStringArg2 := true
	_, err := strconv.Atoi(arg1)
	if err == nil {
		isStringArg1 = false
	}
	_, err = strconv.ParseFloat(arg2, 64)
	if err == nil {
		isStringArg2 = false
	}

	// Use fmt.Sprintf with dynamic format specifiers
	format := "%s + %s = "
	if isStringArg1 || isStringArg2 {
		format += "%s"
	} else {
		format += "%d"
	}

	var result interface{}
	if isStringArg1 || isStringArg2 {
		result = arg1 + arg2
	} else {
		num1, _ := strconv.Atoi(arg1)
		num2, _ := strconv.Atoi(arg2)
		result = num1 + num2
	}
	output := fmt.Sprintf(format, arg1, arg2, result)
	fmt.Println(output)
}
