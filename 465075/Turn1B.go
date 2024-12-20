package main

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
}

func main() {
	user1 := &User{"Alice", 25}
	user2 := (*User)(nil) // nil value
	example1 := fmt.Sprintf("User1: %s, Age: %d", user1.Name, user1.Age)
	example2 := fmt.Sprintf("User2: %v", user2)
	example3 := fmt.Sprintf("User2: %+v", user2)
	fmt.Println(example1)
	fmt.Println(example2)
	fmt.Println(example3)
}
