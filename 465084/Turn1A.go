package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Session represents a session with user data
type Session map[string]interface{}

// Serialize encodes the session data to a binary format
func Serialize(session Session) ([]byte, error) {
	var buf bytes.Buffer // Fixed: use bytes.Buffer instead of io.BytesBuffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(session); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Deserialize decodes the session data from a binary format
func Deserialize(data []byte) (Session, error) {
	var session Session
	decoder := gob.NewDecoder(bytes.NewBuffer(data))

	if err := decoder.Decode(&session); err != nil {
		return nil, err
	}

	return session, nil
}

func main() {
	// Example session data
	sessionData := Session{
		"userId":   "123",
		"username": "johndoe",
		"loggedIn": true,
	}

	// Serialize the session data
	serializedData, err := Serialize(sessionData)
	if err != nil {
		log.Fatalf("Error serializing session data: %v", err)
	}

	// Store the serialized data in a session store (e.g., Redis, database, file)
	// For demonstration, we'll just print it
	log.Printf("Serialized data: %#v\n", serializedData)

	// Deserialize the session data
	deserializedSession, err := Deserialize(serializedData)
	if err != nil {
		log.Fatalf("Error deserializing session data: %v", err)
	}

	// Access the deserialized session data
	log.Printf("Deserialized session: %#v\n", deserializedSession)
}
