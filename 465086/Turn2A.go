package main

import (
	"fmt"
	"sync"
)

// Publisher represents the entity that publishes messages.
type Publisher struct {
	subscribers map[string]map[string]chan string
	lock        sync.RWMutex
}

// NewPublisher creates a new Publisher instance.
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string]map[string]chan string),
	}
}

// Subscribe subscribes a subscriber to a specific topic.
func (p *Publisher) Subscribe(topic string, subscriberID string, subscriber chan string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.subscribers[topic] == nil {
		p.subscribers[topic] = make(map[string]chan string)
	}
	p.subscribers[topic][subscriberID] = subscriber
}

// Publish publishes a message to a specific topic.
func (p *Publisher) Publish(topic string, message string) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if subscribers, ok := p.subscribers[topic]; ok {
		for _, subscriber := range subscribers {
			subscriber <- message
		}
	}
}

// Subscriber represents the entity that subscribes to topics.
func Subscriber(id string, topics []string, publisher *Publisher) {
	subscriber := make(chan string)

	for _, topic := range topics {
		publisher.Subscribe(topic, id, subscriber)
	}

	fmt.Printf("Subscriber %s started for topics: %v\n", id, topics)

	for message := range subscriber {
		fmt.Printf("Subscriber %s received message: %s\n", id, message)
	}
}

func main() {
	publisher := NewPublisher()

	// Start subscribers for different topics
	go Subscriber("S1", []string{"news", "weather"}, publisher)
	go Subscriber("S2", []string{"weather", "finance"}, publisher)
	go Subscriber("S3", []string{"news"}, publisher)
	go Subscriber("S4", []string{"finance"}, publisher)

	// Publish messages to topics
	publisher.Publish("news", "Breaking News: A new update is available!")
	publisher.Publish("weather", "Today's weather forecast: Sunny skies.")
	publisher.Publish("finance", "Stock market closes at a new high.")

	// Wait for subscribers to process messages
	// In a real application, you would use a synchronization mechanism to wait.
	fmt.Scanln()
}
