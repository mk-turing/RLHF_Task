package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func worker(id int, taskCh <-chan string, errCh chan<- error) {
	defer close(errCh)

	for task := range taskCh {
		fmt.Printf("Worker %d processing task: %s\n", id, task)

		// Simulate a delayed task
		time.Sleep(time.Duration(id) * time.Second)

		// Simulate an error
		if task == "errorTask" {
			errCh <- errors.New("task failed")
			return
		}
	}
}

func main() {
	numWorkers := 3
	taskCh := make(chan string, numWorkers*2)
	errChs := make([]chan error, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		errChs[i] = make(chan error, 1)
		go worker(i+1, taskCh, errChs[i])
	}

	// Send tasks
	tasks := []string{"task1", "task2", "errorTask", "task3", "task4"}
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	// Wait for and handle errors
	for _, errCh := range errChs {
		if err := <-errCh; err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}

	fmt.Println("All tasks completed successfully.")
}
