package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Welcome message
	welcome := Message{Message: "Welcome to the chat!"}
	conn.WriteJSON(welcome)

	for {
		var msg Message
		// Read message from client
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
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
		http.ServeFile(w, r, "turn1openapi.yaml")
	})

	http.ListenAndServe(":8080", nil)
}
