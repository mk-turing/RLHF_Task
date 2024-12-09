package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from client
	maxMessageSize = 512

	// Time allowed to write message to peer before closing the connection
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time allowed to write message to the redis channel
	redisWriteWait = 60 * time.Second

	// Broadcast channel name in Redis
	broadcastChannel = "broadcast"
)

var (
	// Mapping of client connections by their unique ID
	clients   = make(map[string]*client)
	muClients sync.RWMutex

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type client struct {
	id            string
	wg            *sync.WaitGroup
	done          chan struct{}
	conn          *websocket.Conn
	send          chan []byte
	redisClient   *redis.Client
	authenticated bool
}

func authenticate(username, password string) bool {
	// Simplified authentication for demonstration
	return username == "test" && password == "test"
}

func serveWS(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	username, password, ok := r.BasicAuth()
	if !ok || !authenticate(username, password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	c := newClient(conn, redisClient, username)

	// Wait for all clients' connections to terminate
	c.wg.Wait()
}

func newClient(conn *websocket.Conn, redisClient *redis.Client, username string) *client {
	var wg sync.WaitGroup
	c := &client{
		id:          username,
		wg:          &wg,
		conn:        conn,
		send:        make(chan []byte),
		redisClient: redisClient,
	}

	wg.Add(2)

	go c.write()
	go c.read()

	return c
}

func (c *client) read() {
	defer c.wg.Done()
	defer func() {
		close(c.done) // Signal the `write` method to stop
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading from client %s: %v", c.id, err)
			}
			return
		}

		// Mark the client as authenticated if they send an authentication message
		if string(message) == "authenticated" {
			c.authenticated = true
			log.Printf("Client %s authenticated", c.id)
			continue
		}

		// Send the message to Redis channel only if client is authenticated
		if !c.authenticated {
			continue
		}

		if err := c.redisClient.RPush(broadcastChannel, message).Err(); err != nil {
			log.Printf("error pushing message to Redis: %v", err)
		}

		c.conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}

func (c *client) write() {
	defer c.wg.Done()

	t := time.NewTicker(pingPeriod)
	defer t.Stop()

	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error writing to client %s: %v", c.id, err)
				return
			}
		case <-t.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("error writing ping to client %s: %v", c.id, err)
				return
			}
		case <-c.done:
			return
		}
	}
}

func broadcast(message []byte) {
	muClients.RLock()
	defer muClients.RUnlock()

	for _, c := range clients {
		select {
		case c.send <- message:
		default:
			log.Printf("client %s queue is full", c.id)
		}
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting to Redis: %v", err)
	}

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			serveWS(w, r, redisClient)
		} else {
			http.ServeFile(w, r, "index.html")
		}
	})))
}
