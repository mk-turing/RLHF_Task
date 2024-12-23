package main

import (
	"fmt"
	"log"
	"os/notify"
)

func main() {
	// Define the directory to watch
	watchDir := "path/to/watch"

	// Create a new watcher
	watcher, err := notify.Watcher.New()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add the directory to watch
	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	defer close(done)

	// Start watching for events
	go func() {
		for {
			select {
			case <-done:
				return
			case event, ok := <-watcher.Events():
				if !ok {
					log.Println("watcher closed")
					return
				}
				// Handle the event
				handleEvent(event)
			case err, ok := <-watcher.Errors():
				if !ok {
					log.Println("watcher closed")
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Main loop to keep the program running
	for {
		select {
		case <-done:
			return
		}
	}
}

func handleEvent(event notify.Event) {
	switch event.Op {
	case notify.Create:
		fmt.Println("File created:", event.Name)
	case notify.Write:
		fmt.Println("File modified:", event.Name)
	case notify.Remove:
		fmt.Println("File removed:", event.Name)
	case notify.Rename:
		fmt.Println("File renamed:", event.Name)
	default:
		fmt.Println("Unknown event:", event.Op)
	}
}
