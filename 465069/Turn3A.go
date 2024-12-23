package main

import (
	"fmt"
)

// formatArgument dynamically formats an argument based on its type
func formatArgument(arg interface{}) string {
	switch arg := arg.(type) {
	case string:
		return fmt.Sprintf("String: %s", arg)
	case int:
		return fmt.Sprintf("Integer: %d", arg)
	case float64:
		return fmt.Sprintf("Float: %.2f", arg)
	default:
		return fmt.Sprintf("Unknown type: %T", arg)
	}
}

func main() {
	// Sample arguments of different types
	var args []interface{} = []interface{}{
		"Hello",
		42,
		3.14159,
		true,
	}

	// Format and print each argument
	for _, arg := range args {
		formatted := formatArgument(arg)
		fmt.Println(formatted)
	}
}
