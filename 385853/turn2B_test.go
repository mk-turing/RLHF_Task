package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelSendReceiv4(t *testing.T) {
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
				defer close(c)
				c <- 1
			},
			receiver: func(c chan int, t *testing.T, wg *sync.WaitGroup) {
				for v := range c {
					require.Contains(t, 1, v)
				}
				wg.Done()
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
