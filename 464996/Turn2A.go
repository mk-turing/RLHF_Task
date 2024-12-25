package main

import (
	"context"
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

// WebSocket handler with proper cleanup and timeout handling
func serveWebSocket(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close() // Ensure connection is closed on exit

	// Context for connection handling
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Ping handling in a separate goroutine
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(1*time.Second)); err != nil {
					log.Println("Ping error:", err)
					cancel()
					return
				}
			}
		}
	}()

	// Client message handling loop
	for {
		select {
		case <-ctx.Done():
			return // Exit when the context is canceled
		default:
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Unexpected client disconnection:", err)
				}
				cancel()
				return
			}
		}
	}
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		serveWebSocket(w, r, ctx)
	})
	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Println("Starting WebSocket server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Test WebSocket connections
	testWebSocketConnections()

	// Shutdown server gracefully
	log.Println("Shutting down WebSocket server")
	server.Close()
}

func testWebSocketConnections() {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			makeWebSocketConnection("ws://localhost:8080/ws")
		}()
	}

	wg.Wait()
	fmt.Printf("Tested %d connections in %v\n", numConnections, time.Since(start))
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
