package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Worker function that performs a task and handles errors with a retry mechanism
func worker(id int, intervalChan <-chan time.Duration, retryChan <-chan time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second) // Default work interval
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			fmt.Printf("Worker %d: Received stop signal, exiting.\n", id)
			return
		case newInterval, ok := <-intervalChan:
			if ok {
				fmt.Printf("Worker %d: Adjusting work interval to %s\n", id, newInterval)
				ticker.Reset(newInterval) // Update the ticker for work interval
			}
		case retryDelay, ok := <-retryChan:
			if ok {
				fmt.Printf("Worker %d: Adjusting retry delay to %s\n", id, retryDelay)
				ticker.Reset(retryDelay) // Update the ticker for retry delay
			}
		case <-ticker.C:
			// Simulate work with a chance of error
			if err := performTask(id); err != nil {
				fmt.Printf("Worker %d: Error occurred: %s. Retrying...\n", id, err)
				time.Sleep(2 * time.Second) // Fixed wait before retrying
				continue                    // Retry immediately after the fixed delay
			}
			fmt.Printf("Worker %d: Task completed successfully.\n", id)
		}
	}
}

// Simulate a task that may fail
func performTask(id int) error {
	if rand.Float32() < 0.3 { // Simulate a 30% chance of failure
		return fmt.Errorf("simulated error in worker %d", id)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	stop := make(chan struct{})
	intervalChan := make(chan time.Duration)
	retryChan := make(chan time.Duration)

	// Start multiple worker goroutines
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		go worker(i, intervalChan, retryChan, stop)
	}

	// Goroutine for user input to adjust intervals and retry delays
	go func() {
		for {
			var input string
			fmt.Print("Enter new work interval (e.g., 2s, 5s), retry delay (e.g., 1s), or 'exit' to stop: ")
			if _, err := fmt.Scanln(&input); err != nil {
				fmt.Println("Error reading input. Please try again.")
				continue
			}
			if input == "exit" {
				close(stop)
				close(intervalChan)
				close(retryChan)
				return
			}
			// Parse the input duration
			newInterval, err := time.ParseDuration(input)
			if err == nil {
				intervalChan <- newInterval // Send the new work interval to workers
				continue
			}
			retryDelay, err := time.ParseDuration(input)
			if err != nil {
				fmt.Println("Invalid format. Please try again.")
				continue
			}
			retryChan <- retryDelay // Send the new retry delay to workers
		}
	}()

	// Main function runs for 30 seconds
	fmt.Println("Main program running. You can adjust work intervals or retry delays, or type 'exit' to stop.")
	time.Sleep(30 * time.Second) // Simulate a running process
	close(stop)
	close(intervalChan)
	close(retryChan)
}
