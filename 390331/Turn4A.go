package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	for i := 0; i < 5; i++ {
		// Simulate producing a message
		message := fmt.Sprintf("Hello from Service A - Message %d", i)

		err = publishMessage(ch, message)
		if err != nil {
			log.Printf("Error publishing message: %v", err)
		}

		time.Sleep(1 * time.Second) // Sleep for demonstration
	}
}

func publishMessage(ch *amqp.Channel, message string) error {
	err := ch.Publish(
		"",           // exchange
		"task_queue", // queue name
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	log.Printf(" [x] Sent %s", message)
	return nil
}
