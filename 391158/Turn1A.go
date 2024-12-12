package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

type Worker struct {
	ID        int
	JobQueue  <-chan Job
	WaitGroup *sync.WaitGroup
}

func (w *Worker) Start() {
	go func() {
		for job := range w.JobQueue {
			fmt.Printf("Worker %d processing job %d\n", w.ID, job.ID)
			// Simulate work
		}
		w.WaitGroup.Done()
	}()
}

func main() {
	const numWorkers = 3
	const numJobs = 10

	jobQueue := make(chan Job, numJobs)
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		worker := Worker{
			ID:        i,
			JobQueue:  jobQueue,
			WaitGroup: &wg,
		}
		worker.Start()
	}

	// Submit jobs
	for i := 1; i <= numJobs; i++ {
		wg.Add(1)
		jobQueue <- Job{ID: i}
	}

	close(jobQueue) // Close the job queue to signal no more jobs

	wg.Wait() // Wait for all jobs to complete
	fmt.Println("All jobs completed.")
}
