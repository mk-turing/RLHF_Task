package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

// Publisher represents the entity that publishes messages to RabbitMQ.
type Publisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewPublisher creates a new Publisher instance.
func NewPublisher() (*Publisher, error) {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{connection: connection, channel: channel}, nil
}

// Close closes the connection and channel.
func (p *Publisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.connection.Close()
}

// Publish publishes a message to a specific topic.
func (p *Publisher) Publish(topic string, message string) error {
	queue, err := p.channel.QueueDeclare(
		topic,
		true,  // Durable: Persist messages
		false, // Exclusive: Private queue
		false, // AutoDelete: Delete queue when no consumers
		false,
		nil, // Arguments: Optional queue arguments
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",         // Exchange
		queue.Name, // Routing key
		false,      // Mandatory: Fail if no queue
		false,      // Immediate: Fail if no consumer
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // TTL: 10 seconds
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("[x] Published %s\n", message)
	return nil
}

// Subscriber represents the entity that consumes messages from RabbitMQ.
func Subscriber(id string, topics []string, numWorkers int) error {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	// Create worker goroutines to process messages
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				msgs, err := channel.Consume(
					"",    // Queue name
					"",    // Consumer name
					true,  // Auto-ack: Automatically acknowledge messages
					false, // Exclusive: Exclusive consumer
					false, // No-local: No local messages
					false, // No-wait: No wait for messages
					nil,   // Arguments: Optional consume arguments
				)
				if err != nil {
					fmt.Println("[x] Consumer closed:", err)
					return
				}

				for d := range msgs {
					fmt.Printf("Subscriber %s processing message: %s\n", id, string(d.Body))
					// Simulate work by adding a random delay
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				}
			}
		}()
	}

	fmt.Printf("Subscriber %s started with %d workers for topics: %v\n", id, numWorkers, topics)

	// Consume messages from all topics
	for _, topic := range topics {
		queue, err := channel.QueueDeclare(
			topic,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		fmt.Printf("[*] Waiting for messages in %s. To exit press CTRL+C\n", queue.Name)
	}

	// Wait for subscribers to process messages
	select {}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	publisher, err := NewPublisher()
	if err != nil {
		fmt.Println("Error creating publisher:", err)
		return
	}
	defer publisher.Close()

	// Start subscribers for different topics with different number of workers
	go Subscriber("S1", []string{"news", "weather"}, 2)
	go Subscriber("S2", []string{"sports", "finance"}, 3)
	go Subscriber("S3", []string{"news", "technology", "entertainment"}, 4)

	// Simulate publishing messages to topics concurrently
	for i := 0; i < 20; i++ {
		topics := []string{"news", "weather", "sports", "finance", "technology", "entertainment"}
		randTopic := topics[rand.Intn(len(topics))]
		message := fmt.Sprintf("Message #%d for topic '%s'", i, randTopic)
		if err := publisher.Publish(randTopic, message); err != nil {
			fmt.Println("Error publishing message:", err)
		}
	}

	// Wait for subscribers to process messages
	// In a real application, you would use a synchronization mechanism to wait.
	fmt.Scanln()
}
