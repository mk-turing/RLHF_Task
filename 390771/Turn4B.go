package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	url          = "https://jsonplaceholder.typicode.com/todos/1"
	numProducers = 5
	numConsumers = 10
	taskQueue    = make(chan string, 100) // Buffered channel to manage tasks
	results      = make([]string, 0)      // Shared data structure for results
	mu           = sync.Mutex{}           // Mutex to guard access to results
)

// Producer function generates tasks and sends them to the task queue
func producer() {
	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer resp.Body.Close()

		// Simulate a task
		task := fmt.Sprintf("Response status: %s", resp.Status)
		taskQueue <- task

		// Introduce a random delay to simulate task generation time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

// Consumer function retrieves tasks from the queue, processes them, and updates shared state
func consumer() {
	for task := range taskQueue {
		// Simulate task processing time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Lock the mutex to modify the shared data structure safely
		mu.Lock()
		results = append(results, task)
		mu.Unlock()
	}
}

func main() {
	fmt.Println("Starting producers and consumers...")

	// Start producer goroutines
	for i := 0; i < numProducers; i++ {
		go producer()
	}

	// Start consumer goroutines
	for i := 0; i < numConsumers; i++ {
		go consumer()
	}

	// Simulate running the system for a fixed period
	time.Sleep(10 * time.Second)

	// Gracefully stop the program
	close(taskQueue)

	fmt.Println("\nProcessing complete. Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Printf("Total Tasks Processed: %d\n", len(results))
}
