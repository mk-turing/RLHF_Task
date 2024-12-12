package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func worker(id int, intervalChan <-chan time.Duration, stop <-chan struct{}) {
	interval := <-intervalChan

	for {
		select {
		case newInterval := <-intervalChan:
			fmt.Printf("Worker %d: Updated interval to %s.\n", id, newInterval)
			interval = newInterval
		case <-stop:
			fmt.Printf("Worker %d received stop signal. Exiting.\n", id)
			return
		default:
			fmt.Printf("Worker %d is working...\n", id)
			time.Sleep(interval) // Simulate work with dynamic interval
		}
	}
}

func main() {
	// Create channels for communication
	stop := make(chan struct{})
	intervalChan := make(chan time.Duration)

	// Start a few worker goroutines with an initial interval
	for i := 1; i <= 3; i++ {
		go worker(i, intervalChan, stop)
	}

	// Initialize the interval channel with a starting value
	intervalChan <- 1 * time.Second

	// Read user input to dynamically change the interval
	fmt.Println("Main task starting... Type 'exit' to stop or a new interval in seconds (e.g., 2):")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "exit":
			fmt.Println("Stopping workers...")
			close(stop)
			close(intervalChan)
			return
		default:
			newInterval, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input. Please enter a valid number or 'exit'.")
				continue
			}
			intervalChan <- time.Duration(newInterval) * time.Second
			fmt.Printf("Updated interval to %d seconds.\n", newInterval)
		}
	}
}
