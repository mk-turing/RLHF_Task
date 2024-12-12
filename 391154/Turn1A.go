package main

import (
	"fmt"
	"time"
)

func worker(id int, timer <-chan time.Time) {
	for t := range timer {
		// Each worker does its job when triggered by the timer
		fmt.Printf("Worker %d triggered at %v\n", id, t)
		// Simulate doing some work
		time.Sleep(1 * time.Second)
	}
}

func main() {
	const numWorkers = 2
	ticker := time.NewTicker(2 * time.Second) // Set a ticker to tick every 2 seconds
	defer ticker.Stop()

	// Create a channel to receive timer events
	timerChan := ticker.C

	// Start worker Goroutines
	for i := 1; i <= numWorkers; i++ {
		go worker(i, timerChan)
	}

	// Run for a duration and then exit
	time.Sleep(10 * time.Second) // Let the program run for 10 seconds
	fmt.Println("Main function ends.")
}
