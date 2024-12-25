package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const numConnections = 10000
const pongWait = 60 * time.Second
const pingPeriod = (pongWait * 9) / 10 // Send pings 90% before pong wait

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Todo: Implement proper origin check
}

// WebSocket handler without defer for cleanup
func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Set pong handler
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start a goroutine to send periodic pings
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pongWait)); err != nil {
					log.Println("WriteControl error:", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Simulate processing
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Println("Unexpected close error:", err)
			}
			return
		}
		// Process the message here
		fmt.Printf("Received message: %s\n", string(message))
	}
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/ws", serveWebSocket)
	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Println("Starting WebSocket server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Test WebSocket connections
	testWebSocketConnections("/ws")

	// Shutdown server gracefully
	log.Println("Shutting down WebSocket server")
	server.Close()
}

func testWebSocketConnections(endpoint string) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			makeWebSocketConnection("ws://localhost:8080" + endpoint)
		}()
	}

	wg.Wait()
	fmt.Printf("Tested %d connections on %s in %v\n", numConnections, endpoint, time.Since(start))
}

func makeWebSocketConnection(url string) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}