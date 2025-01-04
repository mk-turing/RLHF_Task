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

func main() {
	var s MyState

	// Unbuffered channel for immediate synchronization
	immediateChan := make(chan int)
	// Buffered channel with a size of 5 for queued updates
	bufferedChan := make(chan int, 5)

	go func() {
		for {
			select {
			case value, ok := <-immediateChan:
				if !ok {
					return
				}
				s.Lock()
				s.Count += value
				s.Unlock()
			case value, ok := <-bufferedChan:
				if !ok {
					return
				}
				s.Lock()
				s.Count += value
				s.Unlock()
			}
		}
	}()

	// Simulate immediate updates
	go func() {
		for i := 1; i <= 3; i++ {
			immediateChan <- i
		}
		close(immediateChan)
	}()

	// Simulate buffered updates
	go func() {
		for i := 4; i <= 8; i++ {
			bufferedChan <- i
		}
		close(bufferedChan)
	}()

	// Ensure all updates are processed
	time.Sleep(time.Millisecond * 200)

	fmt.Println("Final count:", s.Count) // Output should be Final count: 28
}
