package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type MySession map[string]interface{}

func (m MySession) Serialize() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error serializing session: %w", err)
	}
	return data, nil
}

func (m *MySession) Deserialize(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("error deserializing session: %w", err)
	}
	return nil
}

func main() {
	// Create an instance of MySession
	session := MySession{
		"user_id":  123,
		"username": "guest",
		"email":    "guest@example.com",
	}

	// Serialize the session state to a JSON string
	data, err := session.Serialize()
	if err != nil {
		log.Println("Error serializing session:", err)
		return
	}
	fmt.Println("Serialized session data:", string(data))

	// Let's intentionally change the map structure to demonstrate error handling
	session["new_field"] = "unexpected value"

	// Deserialize the session state from the JSON string
	var deserializedSession MySession
	err = deserializedSession.Deserialize(data)
	if err != nil {
		log.Println("Error deserializing session:", err)
		return
	}
	fmt.Println("Deserialized session data:", deserializedSession)
}
