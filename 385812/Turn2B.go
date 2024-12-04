package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	bufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()
	_, reader, err := conn.NextReader()
	for {
		buf := bufPool.Get().([]byte)
		n, readErr := reader.Read(buf)
		if readErr != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error:", err)
			}
			break
		}

		processPayload(buf[:n])

		// Return the buffer to the pool for recycling
		bufPool.Put(buf)
	}
}

func processPayload(payload []byte) {
	// Process the payload data here
	log.Println("Received payload:", string(payload))
}

func main() {
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
