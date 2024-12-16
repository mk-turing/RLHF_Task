package main

import (
	"fmt"
	"sync"
)

type Event struct {
	ID   int
	Data string
}

type EventHandler struct {
	mu       sync.Mutex
	eventLog []Event
}

func (eh *EventHandler) HandleEvent(event Event) error {
	// Simulate some processing logic that might fail
	if event.Data == "error" {
		return fmt.Errorf("error processing event ID %d", event.ID)
	}

	// Lock the mutex before accessing shared resource
	eh.mu.Lock()
	defer eh.mu.Unlock()

	// Safely append event to the log
	eh.eventLog = append(eh.eventLog, event)
	return nil
}

func main() {
	eventHandler := EventHandler{}
	var wg sync.WaitGroup
	eventChan := make(chan Event)

	// Goroutine to process events
	go func() {
		for event := range eventChan {
			if err := eventHandler.HandleEvent(event); err != nil {
				fmt.Printf("Failed to handle event: %s\n", err)
			}
		}
	}()

	// Simulate sending events
	events := []Event{
		{ID: 1, Data: "normal"},
		{ID: 2, Data: "error"}, // This will trigger an error
		{ID: 3, Data: "normal"},
	}

	// Sending events concurrently
	for _, event := range events {
		wg.Add(1)
		go func(e Event) {
			defer wg.Done()
			eventChan <- e
		}(event)
	}

	// Wait for all events to be sent
	wg.Wait()
	close(eventChan) // Close the channel after sending all events

	// Wait for processing to finish
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Use a small delay to ensure processing completes
		// In a real application, you would want a more robust way to signal completion
		for {
			if len(eventHandler.eventLog) == len(events)-1 { // since one will error out
				break
			}
		}
	}()
	wg.Wait()

	// Display processed events
	fmt.Println("Processed events:", eventHandler.eventLog)
}
