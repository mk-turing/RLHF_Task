package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Function simulating a task that may occasionally fail
func simulateWork() error {
	if rand.Intn(10) == 0 {
		return fmt.Errorf("Work failed randomly")
	}
	return nil
}

// Worker function with retry mechanism
func worker(id int, intervalChan <-chan time.Duration, retryDelayChan <-chan time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second) // Default to 1 second
	defer ticker.Stop()

	retryDelay := time.Second // Default retry delay

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
		case newRetryDelay, ok := <-retryDelayChan:
			if ok {
				retryDelay = newRetryDelay
				fmt.Printf("Worker %d: Adjusting retry delay to %s\n", id, retryDelay)
			}
		case <-ticker.C:
			if err := simulateWork(); err != nil {
				fmt.Printf("Worker %d: Task failed, retrying in %s...\n", id, retryDelay)
				time.Sleep(retryDelay)
			} else {
				fmt.Printf("Worker %d: Performing work...\n", id)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	stop := make(chan struct{})
	intervalChan := make(chan time.Duration)
	retryDelayChan := make(chan time.Duration)

	// Create and start multiple worker goroutines
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		go worker(i, intervalChan, retryDelayChan, stop)
	}

	// Goroutine for user input to adjust intervals and retry delays
	go func() {
		for {
			var input string
			fmt.Print("Enter command (interval <duration>|retryDelay <duration>|exit): ")
			if _, err := fmt.Scanln(&input); err != nil {
				fmt.Println("Error reading input. Please try again.")
				continue
			}

			parts := strings.Split(input, " ")
			switch parts[0] {
			case "interval":
				newInterval, err := time.ParseDuration(parts[1])
				if err != nil {
					fmt.Println("Invalid interval format. Please try again.")
				} else {
					intervalChan <- newInterval
				}
			case "retryDelay":
				newRetryDelay, err := time.ParseDuration(parts[1])
				if err != nil {
					fmt.Println("Invalid retry delay format. Please try again.")
				} else {
					retryDelayChan <- newRetryDelay
				}
			case "exit":
				close(stop)
				close(intervalChan)
				close(retryDelayChan)
				return
			default:
				fmt.Println("Unknown command. Use 'interval <duration>', 'retryDelay <duration>', or 'exit'.")
			}
		}
	}()

	// Main function runs for a while to allow interaction
	fmt.Println("Main program running. You can adjust intervals, retry delays, or type 'exit' to stop.")
	time.Sleep(30 * time.Second) // Simulation of a running process
	close(stop)
	close(intervalChan)
	close(retryDelayChan)
}
