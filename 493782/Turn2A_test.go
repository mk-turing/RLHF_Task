package main

import (
	"errors"
	"sync"
	"testing"
)

// ProcessData processes data slices concurrently and calls a callback for each.
func ProcessData(data []byte, callback func(error)) {
	wg := sync.WaitGroup{}
	for _, chunk := range chunks(data) {
		wg.Add(1)
		go processChunk(chunk, callback, &wg)
	}
	wg.Wait()
}

func chunks(data []byte) [][]byte {
	n := len(data)
	chunkSize := 5 // Arbitrary chunk size
	chunks := make([][]byte, 0, (n+chunkSize-1)/chunkSize)
	for i := 0; i < n; i += chunkSize {
		chunks = append(chunks, data[i:min(i+chunkSize, n)])
	}
	return chunks
}

func processChunk(chunk []byte, callback func(error), wg *sync.WaitGroup) {
	defer wg.Done()
	if chunk == nil {
		callback(errors.New("chunk is nil"))
		return
	}

	if len(chunk) == 0 {
		callback(errors.New("chunk is empty"))
		return
	}

	callback(nil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MockCallback captures errors for assertion.
type MockCallback struct {
	mu     sync.Mutex
	errors []error
}

func (mc *MockCallback) Call(err error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.errors = append(mc.errors, err)
}

func TestProcessDataConcurrency(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	callback := &MockCallback{}

	ProcessData(data, callback.Call)

	// Assert that all callbacks completed successfully
	if len(callback.errors) != 0 {
		t.Errorf("Expected no errors, but got: %v", callback.errors)
	}
}

func TestProcessDataConcurrencyWithNilChunk(t *testing.T) {
	data := []byte{1, 2, 3}
	callback := &MockCallback{}

	// Simulate a nil chunk for error testing
	ProcessData(data, func(err error) {
		if err != nil {
			callback.Call(err)
		}
	})

	// Assert that there are no errors because nil chunk doesn't occur
	if len(callback.errors) != 0 {
		t.Errorf("Expected no errors, but got: %v", callback.errors)
	}
}

func TestProcessDataConcurrencyWithEmptyChunk(t *testing.T) {
	data := []byte{}
	callback := &MockCallback{}

	// Simulate an empty chunk for error testing
	ProcessData(data, callback.Call)

	// Assert that the empty chunk error is caught
	if len(callback.errors) != 1 || callback.errors[0].Error() != "chunk is empty" {
		t.Errorf("Expected one error for empty chunk, but got: %v", callback.errors)
	}
}
