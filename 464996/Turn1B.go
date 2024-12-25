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

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocket handler using defer for cleanup
func serveWebSocketDefer(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Simulate processing
	time.Sleep(time.Millisecond * 100)
	log.Println("Processed connection using defer")
}

// WebSocket handler avoiding defer for cleanup
func serveWebSocketNoDefer(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Simulate processing
	time.Sleep(time.Millisecond * 100)
	conn.Close()
	log.Println("Processed connection without defer")
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/ws/defer", serveWebSocketDefer)
	http.HandleFunc("/ws/nodef", serveWebSocketNoDefer)
	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Println("Starting WebSocket server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Test WebSocket connections
	testWebSocketConnections("/ws/defer")
	testWebSocketConnections("/ws/nodef")

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
	defer conn.Close()

	// Simulate client activity
	time.Sleep(time.Millisecond * 50)
}
