
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
		var ws;
		function connectWebSocket() {
			ws = new WebSocket("ws://localhost:8080/ws");
			ws.onmessage = function(event) {
				var data = JSON.parse(event.data);
				document.getElementById("title").textContent = data.Title;
				document.getElementById("message").textContent = data.Message;
			};
			ws.onclose = function() {
				setTimeout(connectWebSocket, 1000); // Reconnect after 1 second
			};
		}
		connectWebSocket();

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

	client := &Client{
		Conn:        conn,
		Send:        make(chan update),
		lastUpdated: time.Now(),
	}
	addClient(client)
	defer removeClient(client)

	sendUpdate(client, update{Title: "Initial Update", Message: "WebSocket connection established."})

	go handleMessages(client)
}
func handleMessages(client *Client) {
	defer close(client.Send)
	for {
		var message update
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			// Handle connection errors
			log.Println("Error reading message from client:", err)