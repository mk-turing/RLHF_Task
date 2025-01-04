package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type MyState struct {
	Count int
	sync.Mutex
}

type Update struct {
	Value int
	Err   error
}

func main() {
	var s MyState

	updateChan := make(chan Update)

	// Worker function to process updates and errors
	go func() {
		for update := range updateChan {
			s.Lock()
			if update.Err != nil {
				// Handle error asynchronously
				fmt.Println("Error:", update.Err)
			} else {
				s.Count += update.Value
			}
			s.Unlock()
		}
	}()

	// Simulate updates from multiple sources
	for i := 1; i <= 5; i++ {
		go func(i int) {
			// Random delay to introduce variability
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			update := Update{}
			if rand.Intn(2) == 0 {
				// Simulate a successful update
				update.Value = i
			} else {
				// Simulate an error
				update.Err = fmt.Errorf("random error in update %d", i)
			}
			updateChan <- update
		}(i)
	}

	// Close the update channel after all updates have been simulated
	time.Sleep(time.Millisecond * 150)
	close(updateChan)

	fmt.Println("Final count:", s.Count)
}
