package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// simulateWork simulates some work, returning an error if the context is cancelled
func simulateWork(ctx context.Context, id int, delay time.Duration) error {
	select {
	case <-time.After(delay):
		fmt.Printf("Work item %d completed.\n", id)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// throttleWithContext runs work items with throttling, limiting the speed by the delay
func throttleWithContext(maxWorkItems int, workDelay time.Duration) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < maxWorkItems; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := simulateWork(ctx, id, time.Duration(rand.Intn(int(workDelay.Seconds()*1000)))+workDelay); err != nil {
				log.Printf("Work item %d cancelled: %v\n", id, err)
			}
		}(i)
	}

	// Simulate time passing and decide whether to load shed
	time.Sleep(2 * workDelay)
	fmt.Println("System under load, initiating load shedding...")
	cancel() // Cancel all pending work

	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Example usage
	maxWorkItems := 10
	workDelay := time.Second * 1 // Each work item takes 1 second on average
	throttleWithContext(maxWorkItems, workDelay)

	// If you want to test non-context based throttling as a baseline, replace with typical sync/chan patterns for throttling
}
