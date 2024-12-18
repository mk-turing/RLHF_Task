package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	const (
		host    = "localhost"
		port    = "8080"
		timeout = time.Second * 5
	)

	fmt.Println("Connecting to", host, ":", port)

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Could not connect:", err)
	}
	defer conn.Close()

	request := "Hello, World!"
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatal("Failed to write to the connection:", err)
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		log.Fatal("Failed to read from the connection:", err)
	}
	fmt.Println(string(response[:n]))
}
