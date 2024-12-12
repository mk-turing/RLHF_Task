package main

import (
	"fmt"
	"time"
)

func worker(id int, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Printf("Worker %d received stop signal. Exiting.\n", id)
			return
		default:
			fmt.Printf("Worker %d is working...\n", id)
			time.Sleep(time.Second) // Simulate work
		}
	}
}

func main() {
	// Create a stop channel to signal workers to stop
	stop := make(chan struct{})

	// Start a few worker goroutines
	for i := 1; i <= 3; i++ {
		go worker(i, stop)
	}

	// Simulate a main task that runs for 5 seconds
	fmt.Println("Main task starting...")
	time.Sleep(5 * time.Second)
	fmt.Println("Main task completed. Stopping workers...")

	// Close the stop channel to signal workers to exit
	close(stop)
}
