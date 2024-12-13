package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	numProducers = 5
	numConsumers = 3
	numRequests  = 10 // Total requests for producers to generate
	url          = "https://jsonplaceholder.typicode.com/todos/1"
)

type Task struct {
	id int
}

var results []string
var wg sync.WaitGroup
var mu sync.Mutex // Mutex to protect shared access to results

// Producer function to generate tasks
func producer(taskChan chan<- Task) {
	defer wg.Done()
	for i := 0; i < numRequests/numProducers; i++ {
		taskChan <- Task{id: i}
	}
}

// Consumer function to process tasks
func consumer(taskChan <-chan Task) {
	defer wg.Done()
	for task := range taskChan {
		start := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer resp.Body.Close()

		// Record result into shared state with a mutex lock
		mu.Lock()
		results = append(results, fmt.Sprintf("Task ID: %d -> Response status: %s", task.id, resp.Status))
		mu.Unlock()

		// Log latency for task processing
		fmt.Printf("Processed: Task ID: %d, Status: %s, Duration: %v\n", task.id, resp.Status, time.Since(start))
	}
}

func main() {
	taskChan := make(chan Task)

	// Start producer goroutines
	wg.Add(numProducers)
	for i := 0; i < numProducers; i++ {
		go producer(taskChan)
	}

	// Start consumer goroutines
	wg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go consumer(taskChan)
	}

	// Wait for producers to finish
	wg.Wait()
	close(taskChan) // Close the channel after all tasks are produced

	// Wait for consumers to finish, calling WaitGroup without producing any additional tasks
	wg.Wait()

	// Print results
	mu.Lock()
	fmt.Println("Final Results:")
	for _, result := range results {
		fmt.Println(result)
	}
	mu.Unlock()
}
