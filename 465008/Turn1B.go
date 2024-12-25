package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var tmpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Dynamic HTML Update</title>
</head>
<body>
	<h1>{{.Title}}</h1>
	<p>{{.Message}}</p>
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins for simplicity
}

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

	// Send initial update
	sendUpdate(conn, update{Title: "Initial Update", Message: "WebSocket connection established."})

	// Simulate dynamic updates
	for i := 0; i < 5; i++ {
		sendUpdate(conn, update{Title: fmt.Sprintf("Update %d", i+1), Message: "This is update number " + fmt.Sprint(i+1)})
	}
}

func sendUpdate(conn *websocket.Conn, u update) {
	err := conn.WriteJSON(u)
	if err != nil {
		log.Println(err)
	}
}
