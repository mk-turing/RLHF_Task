package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var openapiSpec = `
openapi: 3.0.0
info:
  title: WebSocket API Example
  version: 1.0.0
paths:
  /ws:
    get:
      summary: Connect to WebSocket
      operationId: connectToWebSocket
      responses:
        '200':
          description: WebSocket connection established
components:
  schemas:
    Message:
      type: object
      properties:
        type:
          type: string
        content:
          type: string
  // You can define other components like parameters, request bodies, etc.
`

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for simplicity
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}

		var message struct {
			Type    string `json:"type"`
			Content string `json:"content"`
		}

		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			log.Println("Error unmarshalling WebSocket message:", err)
			continue
		}

		fmt.Printf("Received message: type=%s, content=%s\n", message.Type, message.Content)

		// Example response
		response := map[string]interface{}{
			"type":    "response",
			"content": fmt.Sprintf("Received: %s", message.Content),
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling WebSocket response:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			log.Println("Error writing WebSocket response:", err)
			break
		}
	}
}

func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(openapiSpec))
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/openapi.json", serveOpenAPISpec)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
