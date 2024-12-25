package main

import (
	"encoding/json"
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
	<textarea id="input" rows="4" cols="50"></textarea><br>
	<button onclick="sendMessage()">Send</button>
	<script>
		var ws = new WebSocket("ws://localhost:8080/ws");
		ws.onmessage = function(event) {
			var data = JSON.parse(event.data);
			document.getElementById("title").textContent = data.Title;
			document.getElementById("message").textContent = data.Message;
		};
		function sendMessage() {
			var message = document.getElementById("input").value;
			ws.send(JSON.stringify({Title: "Broadcast", Message: message}));
			document.getElementById("input").value = "";
		}
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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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

	client := &Client{
		Conn:        conn,
		Send:        make(chan update),
		lastUpdated: time.Now(),
	}

	addClient(client)
	defer removeClient(client)

	go handleMessages(client)
}

func handleMessages(client *Client) {
	defer close(client.Send)
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var u update
		if err := json.Unmarshal(message, &u); err != nil {
			log.Println(err)
			continue
		}

		// Check if it's a broadcast message
		if u.Title == "Broadcast" {
			broadcastMessage(u.Message)
		} else {
			// Handle regular updates
			sendUpdate(client, u)
		}
	}
}

func addClient(client *Client) {
	clientsMutex.Lock()