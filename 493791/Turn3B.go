package main
import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Callback func()

type Priority int
const (
	High Priority = iota
	Medium
	Low
)

type CallbackEntry struct {
	callback Callback
	priority Priority
}

type CallbackManager struct {
	callbacks []CallbackEntry
	mu       sync.Mutex
	wg       sync.WaitGroup
	prioQueue []CallbackEntry
}

func (cm *CallbackManager) Add(cb Callback, priority Priority) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, CallbackEntry{cb, priority})
}

// AdjustPriority dynamically changes the priority of an existing callback
func (cm *CallbackManager) AdjustPriority(cb Callback, newPriority Priority) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for i, entry := range cm.callbacks {
		if entry.callback == cb {
			cm.callbacks[i].priority = newPriority
			return
		}
	}
	log.Println("Callback not found to adjust priority.")
}

// Execute executes callbacks based on priority
func (cm *CallbackManager) Execute(numWorkers int) {
	// Initialize priority queue
	cm.prioQueue = cm.callbacks
	for i := len(cm.prioQueue) / 2; i >= 0; i-- {
		cm.siftDown(i)
	}

	cm.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				cm.mu.Lock()
				if len(cm.prioQueue) == 0 {
					cm.mu.Unlock()
					break
				}
				entry := cm.prioQueue[0]
				cm.prioQueue = cm.prioQueue[1:]
				cm.mu.Unlock()
				entry.callback()
			}
			cm.wg.Done()
		}()
	}
	cm.wg.Wait()
}

// siftDown maintains the priority queue property after adding or adjusting a priority
func (cm *CallbackManager) siftDown(i int) {
	n := len(cm.prioQueue)
	for 2*i+1 < n {
		j := 2*i + 1
		if j+1 < n && cm.prioQueue[j].priority < cm.prioQueue[j+1].priority {
			j++
		}
		if cm.prioQueue[i].priority >= cm.prioQueue[j].priority {
			break
		}
		cm.prioQueue[i], cm.prioQueue[j] = cm.prioQueue[j], cm.prioQueue[i]
		i = j
	}
}

func TraceCallback(cb Callback, name string) Callback {
	return func() {
		log.Printf("Starting callback: %s\n", name)
		cb()
		log.Printf("Finished callback: %s\n", name)
	}
}

var raceConditionCount int64

func simulateRaceCondition() {
	atomic.AddInt64(&raceConditionCount, 1)
}

func main() {
	cm := &CallbackManager{}

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 1: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
		simulateRaceCondition()
	}, Medium), "Callback 1")

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 2: Sleeping for 2 seconds...\n")
		time.Sleep(2 * time.Second)
		simulateRaceCondition()