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
		log.Printf("Error serializing session: %v\n", err)
		return nil, err
	}
	return data, nil
}

func (m *MySession) Deserialize(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Printf("Error deserializing session: %v\n", err)
		return err
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
		return
	}
	fmt.Println("Serialized session data:", string(data))

	// Deserialize the session state from the JSON string
	var deserializedSession MySession
	err = deserializedSession.Deserialize(data)
	if err != nil {
		return
	}
	fmt.Println("Deserialized session data:", deserializedSession)
}
