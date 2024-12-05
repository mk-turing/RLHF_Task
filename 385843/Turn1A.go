package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"os"
)

type User struct {
	Name string
	Age  int
}

func validateUser(user User) error {
	if user.Age <= 0 {
		return fmt.Errorf("age should be positive")
	}
	return nil
}

func parseAndValidate(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the entire file into a byte slice
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Parse JSON and validate
	_, err = jsonparser.ArrayEach(fileBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			fmt.Println("Error parsing array element:", err)
			return
		}

		// Assuming the user object is an array element
		user := User{}
		user.Name, _ = jsonparser.GetString(value, "name")
		age, _ := jsonparser.GetInt(value, "age")
		user.Age = int(age)

		// Validate the user data
		err = validateUser(user)
		if err != nil {
			fmt.Println("Validation error:", err)
		} else {
			fmt.Println("User is valid:", user)
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := parseAndValidate("../385843/large_data_file.json"); err != nil {
		fmt.Printf("Error parsing: %s\n", err)
	} else {
		fmt.Println("File parsed and validated successfully")
	}
}
