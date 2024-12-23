package main

import (
	"fmt"
	"log"

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

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Determine the directory path for the event
				for dirPath, _ := range dirActions {
					if event.Name == dirPath || event.Name[len(dirPath):] == "/"+event.Name[len(dirPath):] {
						dirActions[dirPath](event.Name)
						break
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	<-done
}
