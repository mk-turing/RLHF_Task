package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestChannelSendReceive demonstrates how to use parameterized tests for
// different channel send and receive use cases
func TestChannelSendReceive(t *testing.T) {
	tests := []struct {
		name     string
		sender   func(chan int, *sync.WaitGroup)
		receiver func(chan int, *sync.WaitGroup)
		expect   []int
	}{
		{
			name:   "simple send receive",
			sender: func(ch chan int, wg *sync.WaitGroup) { defer wg.Done(); ch <- 1; ch <- 2; close(ch) },
			receiver: func(ch chan int, wg *sync.WaitGroup) {
				defer wg.Done()
				for x := range ch {
					expect(t, x, 2)
				}
			},
			expect: []int{1, 2},
		},
		{
			name:     "with timeout",
			sender:   func(ch chan int, wg *sync.WaitGroup) { defer wg.Done(); ch <- 1; time.Sleep(100 * time.Millisecond) },
			receiver: func(ch chan int, wg *sync.WaitGroup) { defer wg.Done(); expect(t, <-ch, 1) },
			expect:   []int{1},
		},
		{
			name: "select channel send",
			sender: func(ch chan int, wg *sync.WaitGroup) {
				defer wg.Done()
				select {
				case ch <- 1:
				case <-time.After(50 * time.Millisecond):
				}
			},
			receiver: func(ch chan int, wg *sync.WaitGroup) {
				defer wg.Done()
				time.Sleep(100 * time.Millisecond)
				expect(t, <-ch, 1)
			},
			expect: []int{1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := make(chan int)
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go test.sender(ch, wg)
			go test.receiver(ch, wg)
			wg.Wait()
		})
	}
}

func expect(t *testing.T, actual, exp int) {
	require.Equal(t, actual, exp)
}

//func main() {
//}
