package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Simulate a long running operation
func longRunningOperation(ctx context.Context, name string) {
	defer fmt.Println(name, "done")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println(name, "Working")
		}
	}
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	const numGoroutines = 3

	for i := 0; i < numGoroutines; i++ {
		name := fmt.Sprintf("worker %d", i)
		go longRunningOperation(ctx, name)
	}
	//Simulate some work
	time.Sleep(2 * time.Second)
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		fmt.Println("Cancelling operations")
		cancel()
	}
	select {
	case <-ctx.Done():
		fmt.Println("Main finished")
	}
}
