package main

import (
	"fmt"
	"log"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

// Custom error type for handling nil values
type NilError struct {
	Message string
}

func (e *NilError) Error() string {
	return e.Message
}

// Utility function to handle nil values of different types
func handleNil(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.IsNil() {
		var typ string
		switch v.Kind() {
		case reflect.Ptr:
			typ = "pointer"
		case reflect.Slice:
			typ = "slice"
		case reflect.Map:
			typ = "map"
		default:
			typ = "value"
		}
		return &NilError{Message: fmt.Sprintf("%s cannot be nil", typ)}
	}
	return nil
}
func processUser(user *User) (string, error) {
	err := handleNil(user)
	if err != nil {
		return "", err
	}
	// Rest of the function implementation...
	return fmt.Sprintf("User: %s, Age: %d", user.Name, user.Age), nil
}
func processSlice(data []int) ([]int, error) {
	err := handleNil(data)
	if err != nil {
		return nil, err
	}
	// Rest of the function implementation...
	return data, nil
}
func main() {
	user1 := &User{"Alice", 25}
	user2 := (*User)(nil)
	data1 := []int{1, 2, 3}
	data2 := ([]int)(nil)
	result1, err1 := processUser(user1)
	if err1 != nil {
		log.Printf("Error processing user1: %v\n", err1)
	} else {
		fmt.Printf("User1: %s\n", result1)
	}
	result2, err2 := processUser(user2)
	if err2 != nil {
		log.Printf("Error processing user2: %v\n", err2)
	} else {
		fmt.Printf("User2: %s\n", result2)
	}
	result3, err3 := processSlice(data1)
	if err3 != nil {
		log.Printf("Error processing data1: %v\n", err3)
	} else {
		fmt.Printf("Data1: %v\n", result3)
	}
	result4, err4 := processSlice(data2)
	if err4 != nil {
		log.Printf("Error processing data2: %v\n", err4)
	} else {
		fmt.Printf("Data2: %v\n", result4)
	}
}
