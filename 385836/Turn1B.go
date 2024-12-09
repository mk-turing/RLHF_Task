package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nats.io/nats.go"
	"github.com/gorilla/websocket"
)

var (
	upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	nc *nats.Conn
)

func connectNATS() {
	var err error
	nc, err = nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	defer nc.Close()

	msgChan := make(chan []byte)

	err = nc.Subscribe("chat.message", func(m *nats.Msg) {
		msg := string(m.Data)
		msgChan <- []byte(msg)
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}

	go func() {
		for msg := range msgChan {
			fmt.Println("Received message:", string(msg))
			// Broadcast message to all WebSocket clients here
		}
	}()
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}