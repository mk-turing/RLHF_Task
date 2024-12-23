package main

import (
	"encoding/json"
	"fmt"
)

type MySession map[string]interface{}

func (m MySession) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MySession) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
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
		fmt.Println("Error serializing session:", err)
		return
	}
	fmt.Println("Serialized session data:", string(data))

	// Deserialize the session state from the JSON string
	var deserializedSession MySession
	err = deserializedSession.Deserialize(data)
	if err != nil {
		fmt.Println("Error deserializing session:", err)
		return
	}
	fmt.Println("Deserialized session data:", deserializedSession)
}
