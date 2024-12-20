package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Publisher represents the entity that publishes messages asynchronously.
type Publisher struct {
	mu          sync.Mutex
	subscribers map[string][]chan string
}

// NewPublisher creates a new Publisher instance.
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe subscribes a subscriber to a specific topic.
func (p *Publisher) Subscribe(topic string, subscriber chan string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers[topic] = append(p.subscribers[topic], subscriber)
}

// Publish publishes a message to a specific topic asynchronously.
func (p *Publisher) Publish(topic string, message string) {
	go func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		for _, subscriber := range p.subscribers[topic] {
			subscriber <- message
		}
	}()
}

// Subscriber represents the entity that processes messages in parallel.
func Subscriber(id string, topics []string, publisher *Publisher, numWorkers int) {
	// Create a channel to distribute messages to workers
	messageChan := make(chan string, 100)

	// Start worker goroutines to process messages
	for i := 0; i < numWorkers; i++ {
		go func() {
			for message := range messageChan {
				fmt.Printf("Subscriber %s processing message: %s\n", id, message)
				// Simulate work by adding a random delay
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}()
	}

	// Create a buffered channel for each topic to receive messages from the publisher
	subscriberChannels := make(map[string]chan string)
	for _, topic := range topics {
		subscriberChannels[topic] = make(chan string, 100)
	}

	// Subscribe to each topic in the list
	for _, topic := range topics {
		publisher.Subscribe(topic, subscriberChannels[topic])
	}

	fmt.Printf("Subscriber %s started with %d workers for topics: %v\n", id, numWorkers, topics)

	// Handle messages from all topics in separate goroutines
	for topic, channel := range subscriberChannels {
		go func(topic string, channel chan string) {
			for message := range channel {
				messageChan <- message
			}
		}(topic, channel)
	}

	// Wait for all goroutines to exit (this will never happen in this example since we don't close the channels)
	// In a real application, you would use a synchronization mechanism to wait.
	select {}
}

func main() {
	publisher := NewPublisher()

	// Start subscribers for different topics with different number of workers
	go Subscriber("S1", []string{"news", "weather"}, publisher, 2)
	go Subscriber("S2", []string{"sports", "finance"}, publisher, 3)
	go Subscriber("S3", []string{"news", "technology", "entertainment"}, publisher, 4)

	// Simulate publishing messages to topics concurrently
	for i := 0; i < 20; i++ {
		topics := []string{"news", "weather", "sports", "finance", "technology", "entertainment"}
		randTopic := topics[rand.Intn(len(topics))]
		message := fmt.Sprintf("Message #%d for topic '%s'", i, randTopic)
		go publisher.Publish(randTopic, message)
	}

	// Wait for subscribers to process messages
	// In a real application, you would use a synchronization mechanism to wait.
	fmt.Scanln()
}
