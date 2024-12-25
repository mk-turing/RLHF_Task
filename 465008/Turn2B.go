package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"text/template"
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
    <script>
        var ws = new WebSocket("ws://localhost:8080/ws");
        ws.onmessage = function(event) {
            var data = JSON.parse(event.data);
            document.getElementById("title").textContent = data.Title;
            document.getElementById("message").textContent = data.Message;
        };
    </script>
</body>
</html>
`))

type update struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

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

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
	}()

	go handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	for {
		err := conn.WriteJSON(update{Title: fmt.Sprintf("Client %d", len(clients)), Message: "Hello from client!"})
		if err != nil {
			log.Println("Error writing to client:", err)
			return
		}
	}
}
