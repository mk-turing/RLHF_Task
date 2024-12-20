package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message represents a message with a topic and content.
type Message struct {
	Topic   string
	Content string
}

// Publisher represents the entity that publishes messages asynchronously.
type Publisher struct {
	mu           sync.Mutex
	subscribers map[string][]chan Message
	unsentMessages map[string][]Message
}

// NewPublisher creates a new Publisher instance.
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers:   make(map[string][]chan Message),
		unsentMessages: make(map[string][]Message),
	}
}

// Subscribe subscribes a subscriber to a specific topic.
func (p *Publisher) Subscribe(topic string, subscriber chan Message) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers[topic] = append(p.subscribers[topic], subscriber)
}

// Publish publishes a message to a specific topic asynchronously.
func (p *Publisher) Publish(topic string, message string) {
	msg := Message{Topic: topic, Content: message}
	go p.publishMessage(msg)
}

func (p *Publisher) publishMessage(msg Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if there are any subscribers for the topic
	if subscribers, ok := p.subscribers[msg.Topic]; ok {
		// Send the message to all subscribers
		for _, subscriber := range subscribers {
			subscriber <- msg
		}
	} else {
		// Store the unsent message for later delivery
		p.unsentMessages[msg.Topic] = append(p.unsentMessages[msg.Topic], msg)
	}
}

// Subscriber represents the entity that processes messages in parallel.
func Subscriber(id string, topics []string, publisher *Publisher, numWorkers int) {
	// Create a channel to distribute messages to workers
	messageChan := make(chan Message, 100)

	// Start worker goroutines to process messages
	for i := 0; i < numWorkers; i++ {
		go func() {
			for message := range messageChan {
				fmt.Printf("Subscriber %s processing message: %s\n", id, message.Content)
				// Simulate work by adding a random delay
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}()
	}

	// Create a buffered channel for each topic to receive messages from the publisher
	subscriberChannels := make(map[string]chan Message)
	for _, topic := range topics {
		subscriberChannels[topic] = make(chan Message, 100)
	}

	// Subscribe to each topic in the list
	for _, topic := range topics {
		publisher.Subscribe(topic, subscriberChannels[topic])
	}

	fmt.Printf("Subscriber %s started with %d workers for topics: %v\n", id, numWorkers, topics)

	// Handle messages from all topics in separate goroutines
	for topic, channel := range subscriberChannels {
		go func(topic string, channel chan Message) {
			for message := range channel {
				messageChan <- message
			}
		}(topic, channel)
	}

	// Process any unsent messages for the subscribed topics
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, topic := range topics {
		if unsentMessages, ok := p.unsentMessages[topic]; ok {
			for _, msg := range unsentMessages {
				messageChan <- msg
			}
			delete(p.unsentMessages, topic)
		}
	}

	// Wait for all goroutines to exit (this will never happen in this example since we don't close the channels)
	// In a real application, you would use a synchronization mechanism to wait.