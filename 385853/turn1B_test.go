package main

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

// TestChannelSendReceive demonstrates how to use parameterized tests for
// different channel send and receive use cases
func TestChannelSendReceive1(t *testing.T) {
	tests := []struct {
		name     string
		sender   func(chan int, *testing.T, *sync.WaitGroup)
		receiver func(chan int, *testing.T, *sync.WaitGroup)
		expected []int
	}{
		{
			name: "Single sender, single receiver",
			sender: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				wg.Done()
				c <- 1
			},
			receiver: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				require.Equal(t, 1, <-c)
				wg.Done()
			},
			expected: []int{1},
		},
		{
			name: "Multiple senders, single receiver",
			sender: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				defer wg.Done()
				c <- 2
			},
			receiver: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				expected := []int{1, 2, 3}
				for _, exp := range expected {
					require.Equal(t, exp, <-c)
				}
				wg.Done()
			},
			expected: []int{1, 2, 3},
		},
		{
			name: "Single sender, multiple receivers",
			sender: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				c <- 4
				wg.Done()
			},
			receiver: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				require.Equal(t, 4, <-c)
				wg.Done()
			},
			expected: []int{4},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := make(chan int)
			var wg sync.WaitGroup

			wg.Add(1)
			go test.sender(c, t, &wg)

			wg.Add(1)
			go test.receiver(c, t, &wg)

			wg.Wait()
			require.Equal(t, test.expected, make([]int, 0))
		})
	}
}

// tests function
//func main() {
//}
