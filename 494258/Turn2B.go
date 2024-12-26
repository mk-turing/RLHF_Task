package main

import (
	"fmt"
	"sync"
)

var (
	count = 0
)

func processChunk(chunk []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, value := range chunk {
		// Perform some operation on the value
		// In this example, we just count the number of even numbers
		if value%2 == 0 {
			count++
		}
	}
}
func main() {
	const dataSize = 10000000
	data := make([]int, dataSize)
	// Fill the data slice with random values
	// ...

	chunkSize := 1000
	numChunks := dataSize / chunkSize
	wg := sync.WaitGroup{}
	wg.Add(numChunks)
	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		go processChunk(data[start:end], &wg)
	}
	wg.Wait()
	fmt.Println("Count of even numbers:", count)
}
