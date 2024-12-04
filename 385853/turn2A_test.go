package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelSendReceive2(t *testing.T) {
	tests := []struct {
		name     string
		sender   func(chan int, *testing.T, *sync.WaitGroup)
		receiver func(chan int, *testing.T, *sync.WaitGroup)
		expected []int
	}{
		{
			name: "Single sender, single receiver",
			sender: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				defer wg.Done()
				c <- 1
			},
			receiver: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				defer wg.Done()
				received := <-c
				require.Equal(t, 1, received)
			},
			expected: []int{1},
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
		})
	}
}
