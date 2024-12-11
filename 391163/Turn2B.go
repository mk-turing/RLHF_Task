package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Send welcome message
	welcome := Message{Message: "Welcome to the chat!"}
	if err := conn.WriteJSON(welcome); err != nil {
		return
	}

	for {
		var msg Message
		// Read message from client
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}

		// Simulate message processing
		if len(msg.Message) > 255 {
			err := conn.WriteJSON(ErrorResponse{Code: 400, Message: "Invalid message format"})
			if err != nil {
				return
			}
			continue
		}

		// Echo the message back
		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	// Serve OpenAPI documentation
	http.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "turn2bopenapi.yaml")
	})

	http.ListenAndServe(":8080", nil)
}
