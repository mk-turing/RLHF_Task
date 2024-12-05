package main

import (
	"fmt"
)

func main() {
	data := make([]int, 10)
	for i := range data {
		data[i] = i
	}

	// Take snapshot 1 by copying the slice
	snapshot1 := make([]int, len(data))
	copy(snapshot1, data)

	fmt.Println("Original data:", data)
	fmt.Println("Snapshot 1:   ", snapshot1)

	// Modify the original data
	data[3] = 100
	data[7] = 200

	fmt.Println("\nModified original data:", data)
	fmt.Println("Snapshot 1 (remains unchanged):  ", snapshot1)

	// Take snapshot 2 using a pointer (shallow copy)
	snapshot2 := &data

	fmt.Println("\nSnapshot 2 using pointer (shallow copy): ", *snapshot2)

	// Modify the data again, note how snapshot 2 changes too
	data[4] = 300

	fmt.Println("Modified original data after snapshot2:", data)
	fmt.Println("Snapshot 2 (changed due to pointer): ", *snapshot2)
}
