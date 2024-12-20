package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Publisher represents the entity that publishes messages.
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

// Publish publishes a message to a specific topic concurrently.
func (p *Publisher) Publish(topic string, message string) {
	go func(topic string, message string) {
		p.mu.Lock()
		defer p.mu.Unlock()
		for _, subscriber := range p.subscribers[topic] {
			subscriber <- message
		}
	}(topic, message)
}

// Subscriber represents the entity that subscribes to topics.
func Subscriber(id string, topics []string, publisher *Publisher, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a bufferred channel for each topic
	subscriberChannels := make(map[string]chan string)
	for _, topic := range topics {
		subscriberChannels[topic] = make(chan string, 100) // Buffer size of 100 messages
		publisher.Subscribe(topic, subscriberChannels[topic])
	}

	fmt.Printf("Subscriber %s started for topics: %v\n", id, topics)

	// Handle messages for each topic in a separate goroutine
	for topic, channel := range subscriberChannels {
		go func(topic string, channel chan string) {
			for message := range channel {
				fmt.Printf("Subscriber %s received message for topic '%s': %s\n", id, topic, message)
				// Simulate message processing (with a delay)
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}
		}(topic, channel)
	}

	// Wait for all goroutines to exit (this will never happen in this example since we don't close the channels)
	// In a real application, you would use a synchronization mechanism to wait.
	select {}
}

func main() {
	publisher := NewPublisher()
	var wg sync.WaitGroup

	// Start subscribers for different topics
	wg.Add(1)
	go Subscriber("S1", []string{"news", "weather"}, publisher, &wg)
	wg.Add(1)
	go Subscriber("S2", []string{"sports", "finance"}, publisher, &wg)
	wg.Add(1)
	go Subscriber("S3", []string{"news", "technology"}, publisher, &wg)

	// Publish messages to topics concurrently
	go func() {
		time.Sleep(1 * time.Second)
		publisher.Publish("news", "Breaking News: A new update is available!")
		time.Sleep(500 * time.Millisecond)
		publisher.Publish("weather", "Today's weather forecast: Sunny skies.")
		time.Sleep(250 * time.Millisecond)
		publisher.Publish("sports", "Game Updates: Score is tied!")
		time.Sleep(750 * time.Millisecond)
		publisher.Publish("finance", "Stock Market Analysis: Buy on dips!")
		time.Sleep(500 * time.Millisecond)
		publisher.Publish("technology", "New Gadget Release: Check it out!")
	}()

	// Wait for subscribers to finish processing (in a real app, you might stop the subscribers with a signal)
	wg.Wait()
}
