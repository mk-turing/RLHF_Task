package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// This is a simple function to be tested
func addConcurrently(a, b int) (int, error) {
	var sum int

	go func() {
		sum = a + b
		time.Sleep(100 * time.Millisecond)
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		return sum, fmt.Errorf("timeout")
	default:
		return sum, nil
	}
}

func TestAddConcurrently(t *testing.T) {
	runtime.GOMAXPROCS(1)

	tests := []struct {
		a, b int
		want int
	}{
		{a: 1, b: 2, want: 3},
		{a: 10, b: 20, want: 30},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("a=%d,b=%d", test.a, test.b), func(t *testing.T) {
			got, err := addConcurrently(test.a, test.b)
			if err != nil {
				t.Fatalf("addConcurrently(%d, %d) error: %v", test.a, test.b, err)
			}
			if got != test.want {
				t.Errorf("addConcurrently(%d, %d) = %d, want %d", test.a, test.b, got, test.want)
			}
		})
	}
}

func TestDeadlock(t *testing.T) {
	runtime.GOMAXPROCS(1)
	done := make(chan struct{})
	defer close(done)

	go func() {
		<-done // never ends without a close
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		t.Error("Test timed out without deadlock detection")
	default:
		t.Log("Deadlock test ended before timeout")
	}
}
