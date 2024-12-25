package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message struct represents a message to be processed asynchronously.
type Message struct {
	Text   string
	Result string
}

func processMessageAsync(msg *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simulate random processing time
	msg.Result = "Processed: " + msg.Text
}

func main() {
	msgChan := make(chan Message)
	var wg sync.WaitGroup

	go func() {
		for msg := range msgChan {
			wg.Add(1)
			go processMessageAsync(&msg, &wg)
		}
	}()

	// Simulate sending messages asynchronously
	for i := 0; i < 10; i++ {
		msg := Message{Text: fmt.Sprintf("Message %d", i)}
		msgChan <- msg
	}

	// Wait for all messages to be processed
	wg.Wait()

	// Close the channel to indicate completion
	close(msgChan)

	// Simulate sending messages synchronously to verify results
	for msg := range msgChan {
		fmt.Println(msg.Result)
	}
}
