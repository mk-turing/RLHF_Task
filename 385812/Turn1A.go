package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024) // Preallocate a fixed-size buffer
	_, reader, err := conn.NextReader()
	for {
		n, readErr := reader.Read(buf)

		if readErr != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error:", err)
			}
			break
		}

		if n > 0 {
			processPayload(buf[:n])
		}
	}
}

func processPayload(payload []byte) {
	// Process the payload data here
	// For example, you can parse JSON, perform operations, etc.

	// Example: Printing the payload
	log.Println("Received payload:", string(payload))
}

func main() {
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
