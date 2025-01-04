package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type MyState struct {
	Count int
	sync.Mutex
}

type Update struct {
	Value int
	Err   error
}

const (
	numWorkers = 5
	numUpdates = 10000
)

func main() {
	var s MyState
	updateChan := make(chan Update, 100)
	var wg sync.WaitGroup

	wg.Add(numWorkers)
	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		go worker(i, &s, updateChan, &wg)
	}

	// Generate updates and send them to the channel
	go func() {
		for i := 1; i <= numUpdates; i++ {
			update := Update{}
			if rand.Intn(2) == 0 {
				update.Value = i
			} else {
				update.Err = fmt.Errorf("random error in update %d", i)
			}
			updateChan <- update
		}
		close(updateChan)
	}()

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("Final count:", s.Count)
}

func worker(id int, s *MyState, updateChan <-chan Update, wg *sync.WaitGroup) {
	defer wg.Done()
	for update := range updateChan {
		s.Lock()
		if update.Err != nil {
			fmt.Println("Worker", id, ": Error:", update.Err)
		} else {
			s.Count += update.Value
		}
		s.Unlock()
	}
}
