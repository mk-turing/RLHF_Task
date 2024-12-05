package main

import (
	"sync"
)

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]int, 100) // Return a new slice of a fixed size when needed
		},
	}

	defer pool.Put(myFunc(pool.Get())) // Return the slice to the pool after use
}

func myFunc(data interface{}) []int {
	slice := data.([]int)
	// Use the slice here
	return slice
}
