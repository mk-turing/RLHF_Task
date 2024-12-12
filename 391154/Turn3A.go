package main

import (
	"fmt"
	"time"
)

// Worker function that listens for timer interval changes and works accordingly
func worker(id int, intervalChan <-chan time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second) // Default to 1 second
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			fmt.Printf("Worker %d: Received stop signal, exiting.\n", id)
			return
		case newInterval, ok := <-intervalChan:
			if ok {
				fmt.Printf("Worker %d: Adjusting interval to %s\n", id, newInterval)
				ticker.Reset(newInterval) // Update the ticker with the new interval
			}
		case <-ticker.C:
			fmt.Printf("Worker %d: Performing work...\n", id)
		}
	}
}

func main() {
	stop := make(chan struct{})
	intervalChan := make(chan time.Duration)

	// Create and start multiple worker goroutines
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		go worker(i, intervalChan, stop)
	}

	// Goroutine for user input to adjust intervals
	go func() {
		for {
			var input string
			fmt.Print("Enter new interval (e.g., 2s, 5s) or 'exit' to stop: ")
			if _, err := fmt.Scanln(&input); err != nil {
				fmt.Println("Error reading input. Please try again.")
				continue
			}
			if input == "exit" {
				close(stop)
				close(intervalChan)
				return
			}
			// Parse the input duration
			newInterval, err := time.ParseDuration(input)
			if err != nil {
				fmt.Println("Invalid interval format. Please try again.")
				continue
			}
			intervalChan <- newInterval // Send the new interval to workers
		}
	}()

	// Main function runs for 30 seconds
	fmt.Println("Main program running. You can type intervals to adjust workers or type 'exit' to stop.")
	time.Sleep(30 * time.Second) // Simulation of a running process
	close(stop)
	close(intervalChan)
}
