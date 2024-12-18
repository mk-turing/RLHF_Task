package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	clients         = make(map[*websocket.Conn]bool) // Map of connected clients
	broadcastChan   = make(chan string)              // Channel for broadcasting messages
	broadcastMutex  = sync.Mutex{}                   // Mutex to ensure order of messages
	broadcastQueue  = []string{}                     // Queue to hold messages before sending
	broadcastTicker = time.NewTicker(50 * time.Millisecond)
)

type client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

// WebSocket connection handler
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client
	c := &client{conn: conn, send: make(chan string, 256)}

	// Register the client
	clients[conn] = true
	log.Printf("New client connected. Total clients: %d\n", len(clients))

	// Start reading and writing messages concurrently
	go c.readPump()
	go c.writePump()
}

// Function to handle broadcasting messages to all connected clients
func handleBroadcast() {
	for {
		message := <-broadcastChan

		// Queue the message to ensure order of delivery
		broadcastMutex.Lock()
		broadcastQueue = append(broadcastQueue, message)
		broadcastMutex.Unlock()

		// Send messages from the queue to clients
		select {
		case <-broadcastTicker.C:
			broadcastMutex.Lock()
			for len(broadcastQueue) > 0 {
				//msg := broadcastQueue[0]
				broadcastQueue = broadcastQueue[1:]
				broadcastMutex.Unlock()

				// Send the message to all connected clients
				//for client := range clients {
				//	//select {
				//	//case client.send <- msg:
				//	//default:
				//	//	// Client's send channel is full, close the connection
				//	//	client.conn.Close()
				//	//	delete(clients, client.conn)
				//	//}
				//}

				// Lock the mutex back for the next message
				broadcastMutex.Lock()
			}
			broadcastMutex.Unlock()
		default:
		}
	}
}

// readPump pumps messages from the websocket connection to the broadcast channel.
func (c *client) readPump() {
	defer func() {
		c.conn.Close()
		delete(clients, c.conn)
		log.Printf("Client disconnected. Total clients: %d\n", len(clients))
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		broadcastChan <- string(message)
	}
}

// writePump pumps messages from the client's send channel to the websocket connection.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The send channel has been closed, close the websocket connection
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte(<-c.send))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func main() {
	// Start the broadcast handler
	go handleBroadcast()

	// Serve the WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}
