package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Allow all CPU cores to be used

	data := make([]int, 10000000) // Large dataset
	for i := range data {
		data[i] = i
	}

	var unoptimizedTime time.Duration
	var optimizedTime time.Duration

	// Unoptimized range loop
	start := time.Now()
	for range data {
	}
	unoptimizedTime = time.Since(start)

	// Optimized loop with stored length
	start = time.Now()
	length := len(data)
	for i := 0; i < length; i++ {
	}
	optimizedTime = time.Since(start)

	fmt.Printf("Unoptimized Time: %s, Optimized Time: %s\n", unoptimizedTime, optimizedTime)
}
