package main

import (
	"fmt"
)

// Publisher represents the entity that publishes messages.
type Publisher struct {
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
	p.subscribers[topic] = append(p.subscribers[topic], subscriber)
}

// Publish publishes a message to a specific topic.
func (p *Publisher) Publish(topic string, message string) {
	for _, subscriber := range p.subscribers[topic] {
		subscriber <- message
	}
}

// Subscriber represents the entity that subscribes to topics.
func Subscriber(id string, topic string, publisher *Publisher) {
	subscriber := make(chan string)
	publisher.Subscribe(topic, subscriber)

	fmt.Printf("Subscriber %s started for topic %s\n", id, topic)

	for message := range subscriber {
		fmt.Printf("Subscriber %s received message: %s\n", id, message)
	}
}

func main() {
	publisher := NewPublisher()

	// Start subscribers for different topics
	go Subscriber("S1", "news", publisher)
	go Subscriber("S2", "weather", publisher)
	go Subscriber("S3", "news", publisher)
	go Subscriber("S4", "weather", publisher)

	// Publish messages to topics
	publisher.Publish("news", "Breaking News: A new update is available!")
	publisher.Publish("weather", "Today's weather forecast: Sunny skies.")

	// Wait for subscribers to process messages
	// In a real application, you would use a synchronization mechanism to wait.
	fmt.Scanln()
}
