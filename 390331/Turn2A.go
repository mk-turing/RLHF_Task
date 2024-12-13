package main

import (
	"fmt"
	"sync"
	"time"
)

// Custom error type for demonstration
type TaskError struct {
	TaskID int
	Err    error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("Error in task %d: %v", e.TaskID, e.Err)
}

// Simulated task function
func runTask(id int, ch chan<- *TaskError) {
	// Simulating some work
	time.Sleep(time.Millisecond * time.Duration(100+id*20))

	// Simulating an error for certain tasks
	if id%3 == 0 { // Let's say every third task fails
		ch <- &TaskError{TaskID: id, Err: fmt.Errorf("simulated error")}
		return
	}

	// If successful, send a nil error
	ch <- nil
}

func main() {
	const numTasks = 10
	ch := make(chan *TaskError, numTasks) // Buffered channel for error communication
	var wg sync.WaitGroup

	// Start tasks
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runTask(id, ch)
		}(i)
	}

	// Close the channel after all tasks are done
	go func() {
		wg.Wait() // Wait for all goroutines to finish
		close(ch) // Close channel when done
	}()

	// Collect results and handle errors
	for err := range ch {
		if err != nil {
			fmt.Println(err) // Handle errors here (log, collect, etc.)
		} else {
			fmt.Printf("Task completed successfully\n")
		}
	}

	fmt.Println("All tasks processed.")
}
