package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	clients = make(map[string]*client)
	mu      sync.Mutex

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Basic authentication credentials
	users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}
)

type client struct {
	wg          *sync.WaitGroup
	conn        *websocket.Conn
	send        chan []byte
	redisClient *redis.Client
	id          string
	username    string
}

func (c *client) read() {
	defer c.wg.Done()
	defer c.conn.Close()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading: %v", err)
			}
			return
		}

		c.broadcast(message)

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
				log.Printf("error writing: %v", err)
				return
			}
		case <-t.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("error writing ping: %v", err)
				return
			}
		case <-c.conn.CloseNotify():
			return
		}
	}
}

func (c *client) broadcast(message []byte) {
	// Publish the message to the Redis broadcast channel
	if err := c.redisClient.RPush(broadcastChannel, message).Err(); err != nil {
		log.Printf("error pushing message to Redis: %v", err)
		return
	}

	// Send the message to all clients
	mu.Lock()
	defer mu.Unlock()

	for _, cli := range clients {
		if cli.username != c.username {
			select {
			case cli.send <- message:
			default:
				log.Println("client send buffer full")
			}
		}
	}
}

func (c *client) cleanUp() {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.id)
	log.Printf("Client %s disconnected", c.username)
}

func serveWS(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	if actualPassword, exists := users[username]; !exists || actualPassword != password {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	c := newClient(conn, redisClient, username)

	mu.Lock()
	clients[c.id] = c
	mu.Unlock()

	go c.write()
	go c.read()

	// Wait for all clients connections to terminate
	c.wg.Wait()
	c.cleanUp()
}

func newClient(conn *websocket.Conn, redisClient *redis.Client, username string) *client {
	var wg sync.WaitGroup
	c := &client{
		wg:          &wg,
		conn:        conn,
		send:        make(chan []byte),
		redisClient: redisClient,
		id:          generateClientID(),
		username:    username,
	}

	wg.Add(2)

	return c
}

func generateClientID() string {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	return hex.EncodeToString(hash.Sum(nil))
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
