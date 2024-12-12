package main

import (
	"fmt"
	"time"
)

func worker(id int, intervalChan <-chan time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second) // Default interval
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			fmt.Printf("Worker %d received stop signal. Exiting.\n", id)
			return
		case newInterval, ok := <-intervalChan:
			if ok {
				fmt.Printf("Worker %d received new interval: %s\n", id, newInterval)
				ticker.Reset(newInterval)
			}
		case <-ticker.C:
			fmt.Printf("Worker %d is working...\n", id)
		}
	}
}

func main() {
	stop := make(chan struct{})
	intervalChan := make(chan time.Duration)

	// Start worker goroutines
	for i := 1; i <= 3; i++ {
		go worker(i, intervalChan, stop)
	}

	// Allow dynamic adjustment of intervals
	go func() {
		for {
			var input string
			fmt.Print("Enter new interval (e.g., 2s, 5s) or 'exit' to stop: ")
			fmt.Scanln(&input)
			if input == "exit" {
				close(stop)
				close(intervalChan)
				return
			}
			newInterval, err := time.ParseDuration(input)
			if err != nil {
				fmt.Println("Invalid interval format. Please try again.")
				continue
			}
			intervalChan <- newInterval
		}
	}()

	// Simulate a main task that runs for a certain period
	fmt.Println("Main task running. Type intervals to adjust workers or 'exit' to stop.")
	time.Sleep(30 * time.Second) // Run for 30 seconds
	close(stop)
	close(intervalChan)
}
