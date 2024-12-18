package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	host              = "localhost"
	port              = "8080"
	timeoutStr        = "TIMEOUT_SECONDS"
	initialTimeout, _ = strconv.Atoi(os.Getenv(timeoutStr))
	minTimeout        = time.Second * 1
	maxTimeout        = time.Second * 10
	latency           time.Duration
	mu                sync.Mutex
)

func main() {
	fmt.Println("Connecting to", host, ":", port)

	for {
		dialer := &net.Dialer{Timeout: time.Second * time.Duration(getTimeout())}
		conn, err := dialer.Dial("tcp", host+":"+port)
		if err != nil {
			log.Println("Could not connect:", err)
			time.Sleep(time.Second * time.Duration(getBackoff()))
			continue
		}
		defer conn.Close()

		start := time.Now()
		request := "Hello, World!"
		_, err = conn.Write([]byte(request))
		if err != nil {
			log.Println("Failed to write to the connection:", err)
			time.Sleep(time.Second * time.Duration(getBackoff()))
			continue
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			log.Println("Failed to read from the connection:", err)
			time.Sleep(time.Second * time.Duration(getBackoff()))
			continue
		}

		end := time.Now()
		latency = end.Sub(start)
		updateTimeout(latency)

		fmt.Println(string(response[:n]))
	}
}

func getTimeout() int {
	mu.Lock()
	defer mu.Unlock()
	if latency > maxTimeout {
		return int(maxTimeout)
	}
	if latency < minTimeout {
		return int(minTimeout)
	}
	return int(latency.Seconds()) + 1
}

func updateTimeout(l time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	latency = l
}

func getBackoff() int {
	backoff := 1
	if latency > maxTimeout {
		backoff = 16 // Exponential backoff with a maximum of 16 seconds
	}
	return backoff
}
