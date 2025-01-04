package main

import (
	"fmt"
	"sync"
	"time"
)

type MyState struct {
	Count int
	sync.Mutex
}

type StateUpdate struct {
	Type  string // e.g., "add" or "error"
	Value int    // value to add or error message
}

func main() {
	var s MyState
	updates := make(chan StateUpdate)
	errors := make(chan StateUpdate)

	go func() {
		for update := range updates {
			switch update.Type {
			case "add":
				s.Lock()
				s.Count += update.Value
				s.Unlock()
			default:
				fmt.Println("Invalid update type:", update.Type)
			}
		}
	}()

	go func() {
		for errorUpdate := range errors {
			fmt.Println("Error:", errorUpdate.Value)
		}
	}()

	// Simulate successful state updates
	go func() {
		for i := 1; i <= 3; i++ {
			updates <- StateUpdate{"add", i}
		}
		updates <- StateUpdate{"error", 100} // Simulate an error
		close(updates)
	}()

	// Simulate more successful state updates
	go func() {
		for i := 4; i <= 8; i++ {
			updates <- StateUpdate{"add", i}
		}
		updates <- StateUpdate{"error", 200} // Simulate another error
		close(updates)
	}()

	// Ensure all updates and errors are processed
	time.Sleep(time.Millisecond * 200)

	fmt.Println("Final count:", s.Count) // Output should be Final count: 28
}
