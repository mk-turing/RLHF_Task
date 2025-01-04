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

func incrementCounter(wg *sync.WaitGroup, msgChan chan<- int, doneChan chan bool) {
	defer wg.Done()
	for {
		select {
		case <-doneChan:
			return
		default:
			msgChan <- 1
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
	}
}

func decrementCounter(wg *sync.WaitGroup, msgChan chan<- int, doneChan chan bool) {
	defer wg.Done()
	for {
		select {
		case <-doneChan:
			return
		default:
			msgChan <- -1
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
	}
}

func updateState(wg *sync.WaitGroup, msgChan <-chan int, s *MyState) {
	defer wg.Done()
	for msg := range msgChan {
		s.Lock()
		s.Count += msg
		s.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	// Create an unbuffered channel for state updates
	msgChan := make(chan int)
	// Create a buffered channel for control signals
	doneChan := make(chan bool, 2)

	var s MyState
	wg.Add(4)
	// Start multiple incrementing and decrementing goroutines
	go incrementCounter(&wg, msgChan, doneChan)
	go incrementCounter(&wg, msgChan, doneChan)
	go decrementCounter(&wg, msgChan, doneChan)
	go decrementCounter(&wg, msgChan, doneChan)

	go updateState(&wg, msgChan, &s)
	time.Sleep(2 * time.Second)
	// Signal the goroutines to exit
	close(doneChan)
	wg.Wait()

	fmt.Println("Final count:", s.Count) // Output should be Final count: 0 (with buffered channel) or close to 0 (with unbuffered channel)
}
