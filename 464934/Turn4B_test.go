package main

import (
	"sync"
	"testing"
)

func TestAsyncProcessing(t *testing.T) {
	msgChan := make(chan Message)
	var wg sync.WaitGroup

	// Test multiple messages
	for _, testMsg := range []string{"Test 1", "Test 2", "Test 3"} {
		wg.Add(1)
		go processMessageAsync(&Message{Text: testMsg}, &wg)
	}

	// Wait for all messages to be processed
	wg.Wait()

	// Close the channel once processing is done
	close(msgChan)

	// Verify results using range to read from the channel
	for msg := range msgChan {
		expectedResult := "Processed: " + msg.Text
		if msg.Result != expectedResult {
			t.Errorf("Expected %q, got %q", expectedResult, msg.Result)
		}
	}
}
