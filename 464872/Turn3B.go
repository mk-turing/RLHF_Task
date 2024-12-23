package main

import (
	"fmt"
	"log"
	"sync"

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

	// Buffered channel to receive events
	events := make(chan fsnotify.Event, 100)

	// WaitGroup to ensure all goroutines finish
	var wg sync.WaitGroup

	// Goroutine to handle events concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			for dirPath, _ := range dirActions {
				if event.Name == dirPath || event.Name[len(dirPath):] == "/"+event.Name[len(dirPath):] {
					dirActions[dirPath](event.Name)
					break
				}
			}
		}
	}()

	// Start watching events
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				events <- event // Send event to the buffered channel
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
