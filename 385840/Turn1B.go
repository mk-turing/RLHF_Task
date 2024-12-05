package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Data is a struct to hold our shared data.
type Data struct {
	sync.Mutex
	count int
}

// increment increments the count in the Data struct with a lock.
func increment(ctx context.Context, d *Data) {
	select {
	case <-ctx.Done():
		fmt.Println("Increment goroutine cancelled.")
		return
	default:
	}

	d.Lock()
	defer d.Unlock()
	d.count++

	// Simulate some work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sharedData := &Data{count: 0}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	const numGoroutines = 10

	// Create multiple goroutines to increment the shared data concurrently.
	for i := 0; i < numGoroutines; i++ {
		go increment(ctx, sharedData)
	}

	// Simulate some main thread work
	time.Sleep(time.Second)

	fmt.Println("Final Count:", sharedData.count)
}
