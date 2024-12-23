package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Define a type for actions
type action func(string)

// Create a map that maps directory paths to their corresponding actions
var dirActions map[string]action

// Initialize the action map
func init() {
	dirActions = make(map[string]action)
	dirActions["dir"] = handleDir1Event
	dirActions["dir1"] = handleDir2Event
	// Add more directories and their corresponding actions as needed
}

// Example action function for directory 1
func handleDir1Event(path string) {
	fmt.Printf("Event detected in directory 1: %s\n", path)
}

// Example action function for directory 2
func handleDir2Event(path string) {
	fmt.Printf("Event detected in directory 2: %s\n", path)
}

var eventCount uint64 = 0                        // Counter for event processing
var eventQueue = make(chan fsnotify.Event, 1000) // Queue for events
var workerCount = 3                              // Number of worker threads
var wg sync.WaitGroup

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add directories to watch
	for dirPath := range dirActions {
		err = watcher.Add(dirPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Watching directory: %s\n", dirPath)
	}

	// Start worker threads
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Determine the directory path for the event
				for dirPath, _ := range dirActions {
					if event.Name == dirPath || event.Name[len(dirPath):] == "/"+event.Name[len(dirPath):] {
						eventQueue <- event
						break
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// Simulate a duration for the application to run
	time.Sleep(10 * time.Second)

	// Shut down worker threads
	close(eventQueue)
	wg.Wait()
}

func worker() {
	for event := range eventQueue {
		atomic.AddUint64(&eventCount, 1)
		dirPath := ""
		for dirPath, _ := range dirActions {
			if event.Name == dirPath || event.Name[len(dirPath):] == "/"+event.Name[len(dirPath):] {
				break
			}
		}

		if dirPath == "" {
			continue // Event did not match any watched directory
		}

		dirActions[dirPath](event.Name)
		time.Sleep(10 * time.Millisecond) // Simulate work time
	}
	wg.Done()
}
