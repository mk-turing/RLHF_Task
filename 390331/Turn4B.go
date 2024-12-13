package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pika/pika"
)

type Message struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Attempts int    `json:"attempts"`
}

func PublishMessage(ctx context.Context, connection *pika.Connection, message Message) error {
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	err = channel.QueueDeclare(
		"my_queue", // name
		false,      // durable
		false,      // exclusive
		false,      // auto_delete
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = channel.BasicPublish(
		"",         // exchange
		"my_queue", // routing key
		false,      // mandatory
		false,      // immediate
		pika.BasicProperties{
			DeliveryMode: pika.PersistentDeliveryMode,
		},
		msgBody,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	connection, err := pika.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer connection.Close()

	msg := Message{
		ID:       "123",
		Content:  "Hello, Consumer!",
		Attempts: 0,
	}

	for {
		err = PublishMessage(ctx, connection, msg)
		if err == nil {
			break
		}

		log.Printf("Failed to publish message, retrying: %v", err)
		time.Sleep(time.Second) // Simple exponential backoff
	}
}
