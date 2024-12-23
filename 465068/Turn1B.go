package main

import "fmt"

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Result[T any] struct {
	Value T
	Error error
}

func SomeFunction[T Number](input int) Result[T] {
	var result Result[T]

	if input == 0 {
		result.Error = fmt.Errorf("Input cannot be zero")
	} else {
		result.Value = T(input * input)
	}
	return result
}

func main() {
	result := SomeFunction[int](5)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	} else {
		fmt.Println("Result:", result.Value)
	}

	result = SomeFunction[int](0)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	} else {
		fmt.Println("Result:", result.Value)
	}
}
