package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Resource struct {
	limit   int
	current int
	mutex   sync.Mutex
}

// Task function simulates a workload requiring certain resources.
func task(ctx context.Context, id int, resourceType string) {
	select {
	case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
		fmt.Printf("Task %d completed (Resource: %s)\n", id, resourceType)
	case <-ctx.Done():
		fmt.Printf("Task %d cancelled (Resource: %s)\n", id, resourceType)
	}
}

// Worker function processes tasks, using a semaphore for each resource type.
func worker(ctx context.Context, wg *sync.WaitGroup, id int, resource *Resource, resourceType string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resource.mutex.Lock()
			if resource.current < resource.limit {
				resource.current++
				resource.mutex.Unlock()
				task(ctx, id, resourceType)
				resource.mutex.Lock()
				resource.current--
			} else {
				resource.mutex.Unlock()
				time.Sleep(50 * time.Millisecond) // Wait before checking again
			}
		}
	}
}

// SimulateLoad monitors and adjusts the load dynamically for each resource.
func simulateLoad(ctx context.Context, resource *Resource, adjustmentInterval time.Duration, loadCh chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			loadValue := rand.Intn(5) + 1 // Simulated load between 1 and 5
			resource.mutex.Lock()
			resource.limit = 5 - loadValue // Adjust limit based on load
			if resource.limit < 1 {
				resource.limit = 1
			}
			fmt.Printf("Adjusting %s limit to %d based on load %d\n", resourceType(resource), resource.limit, loadValue)
			loadCh <- loadValue
			resource.mutex.Unlock()
			time.Sleep(adjustmentInterval)
		}
	}
}

func resourceType(res *Resource) string {
	switch res {
	case cpuResource:
		return "CPU"
	case memoryResource:
		return "Memory"
	case networkResource:
		return "Network"
	default:
		return "Unknown"
	}
}

var (
	cpuResource     = &Resource{limit: 3}
	memoryResource  = &Resource{limit: 3}
	networkResource = &Resource{limit: 3}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const maxWorkers = 10
	const totalTasks = 20
	loadCh := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Starting worker goroutines for different resources
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		resourceType := rand.Intn(3) // Randomly prioritize CPU, Memory, or Network
		switch resourceType {
		case 0:
			go worker(ctx, &wg, i, cpuResource, "CPU")
		case 1:
			go worker(ctx, &wg, i, memoryResource, "Memory")
		case 2:
			go worker(ctx, &wg, i, networkResource, "Network")
		}
	}

	// Simulate resource load adjustments
	go simulateLoad(ctx, cpuResource, 2*time.Second, loadCh)
	go simulateLoad(ctx, memoryResource, 2*time.Second, loadCh)
	go simulateLoad(ctx, networkResource, 2*time.Second, loadCh)

	// Sending tasks
	for i := 0; i < totalTasks; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			// Each task is sent to a random resource type
			resourceRequest := rand.Intn(3) // Randomly select resource type for the task
			switch resourceRequest {
			case 0:
				cpuResource.mutex.Lock()
				if cpuResource.current < cpuResource.limit {
					cpuResource.current++
					cpuResource.mutex.Unlock()
					go task(ctx, i, "CPU")
				} else {
					cpuResource.mutex.Unlock()
					fmt.Printf("Task %d dropped due to CPU load\n", i)
				}
			case 1:
				memoryResource.mutex.Lock()
				if memoryResource.current < memoryResource.limit {
					memoryResource.current++
					memoryResource.mutex.Unlock()
					go task(ctx, i, "Memory")
				} else {
					memoryResource.mutex.Unlock()
					fmt.Printf("Task %d dropped due to Memory load\n", i)
				}
			case 2:
				networkResource.mutex.Lock()
				if networkResource.current < networkResource.limit {
					networkResource.current++
					networkResource.mutex.Unlock()
					go task(ctx, i, "Network")
				} else {
					networkResource.mutex.Unlock()
					fmt.Printf("Task %d dropped due to Network load\n", i)
				}
			}
		}
		time.Sleep(100 * time.Millisecond) // Simulate incoming tasks
	}

	// Wait for all workers to finish
	wg.Wait()
}
