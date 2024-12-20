package main

import "fmt"

// Message is the structure of a message that can be published and subscribed to.
type Message struct {
	Topic   string
	Content string
}

// Subscribe function handles subscriber registration for a specific topic.
var subscriptions map[string][]func(*Message)

func Subscribe(topic string, subscriber func(*Message)) {
	if subscriptions == nil {
		subscriptions = make(map[string][]func(*Message))
	}
	if subscribers, exists := subscriptions[topic]; exists {
		subscriptions[topic] = append(subscribers, subscriber)
	} else {
		subscriptions[topic] = []func(*Message){subscriber}
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

// Publish function sends a message to all subscribers of the specified topic.
func Publish(message *Message) {
	if subscribers, exists := subscriptions[message.Topic]; exists {
		for _, subscriber := range subscribers {
			subscriber(message)
		}
		fmt.Printf("Published message: %s to topic: %s\n", message.Content, message.Topic)
	} else {
		fmt.Printf("No subscribers for topic: %s\n", message.Topic)
	}
}

func main() {
	// Sample subscriber function
	subscriber1 := func(message *Message) {
		fmt.Printf("Subscriber 1 received message: %s on topic: %s\n", message.Content, message.Topic)
	}
	subscriber2 := func(message *Message) {
		fmt.Printf("Subscriber 2 received message: %s on topic: %s\n", message.Content, message.Topic)
	}

	// Register subscribers
	Subscribe("weather", subscriber1)
	Subscribe("weather", subscriber2)
	Subscribe("stocks", subscriber1)

	// Publish messages
	Publish(&Message{Topic: "weather", Content: "It's sunny today."})
	Publish(&Message{Topic: "stocks", Content: "Apple's stock increased."})
	Publish(&Message{Topic: "news", Content: "Breaking news!"}) // No subscribers for this topic
}
