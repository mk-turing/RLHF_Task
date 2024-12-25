package main

import (
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Dynamic HTML Update</title>
</head>
<body>
	<h1 id="title">{{.Title}}</h1>
	<p id="message">{{.Message}}</p>
	<input type="text" id="message-input" placeholder="Enter your message">
	<button id="send-button">Send</button>
	<script>
		var ws = new WebSocket("ws://localhost:8080/ws");
		ws.onmessage = function(event) {
			var data = JSON.parse(event.data);
			document.getElementById("title").textContent = data.Title;
			document.getElementById("message").textContent = data.Message;
		};

		document.getElementById("send-button").addEventListener("click", function() {
			var message = document.getElementById("message-input").value;
			if (message.length > 0) {
				ws.send(JSON.stringify({
					Title: "Message Received",
					Message: "From: " + document.location.href + "\nMessage: " + message
				}));
				document.getElementById("message-input").value = "";
			}
		});
	</script>
</body>
</html>
`))

type update struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins for simplicity
}

// Client represents a connected WebSocket client.
type Client struct {
	Conn        *websocket.Conn
	Send        chan update
	lastUpdated time.Time
}

var clients = make([]*Client, 0)
var clientsMutex sync.RWMutex

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)

	go sendUpdates()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, update{Title: "Welcome", Message: "Connecting to WebSocket..."})
	if err != nil {
		log.Println(err)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Create a new client struct to manage this connection
	client := &Client{
		Conn:        conn,
		Send:        make(chan update),
		lastUpdated: time.Now(),
	}

	// Add client to the list
	addClient(client)
	defer removeClient(client)

	// Send initial update
	sendUpdate(client, update{Title: "Initial Update", Message: "WebSocket connection established."})

	go handleMessages(client)
}

func handleMessages(client *Client) {
	defer close(client.Send)
	for {
		// Correctly handle reading messages from the WebSocket connection
		var message update
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			// If an error occurs (e.g., client disconnects), we log it and stop
			log.Println("Error reading message from client:", err)
			return
		}

		log.Printf("Received message from client: %v\n", message.Message)
		broadcastMessage(message)
	}
}

// Adds a new client to the list
func addClient(client *Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients = append(clients, client)
	log.Printf("New client connected: %v", client.Conn.RemoteAddr())
}

// Removes a client from the list
func removeClient(client *Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for i, c := range clients {
		if c == client {
			log.Printf("Client disconnected: %v", c.Conn.RemoteAddr())
			clients = append(clients[:i], clients[i+1:]...)
			return
		}
	}
}

// Sends an update to a specific client
func sendUpdate(client *Client, u update) {
	client.lastUpdated = time.Now()
	client.Send <- u
}

// Sends an update to all connected clients
func broadcastMessage(update update) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	for _, client := range clients {
		client.Send <- update
	}
}

// Sends periodic updates to all connected clients
func sendUpdates() {
	for {
		time.Sleep(1 * time.Second)
		clientsMutex.RLock()
		for _, client := range clients {
			if client.lastUpdated.Add(10 * time.Second).Before(time.Now()) {
				sendUpdate(client, update{Title: "Inactivity", Message: "This client hasn't received any updates for 10 seconds."})
			}
		}
		clientsMutex.RUnlock()
	}
}
