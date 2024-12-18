package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	host             = "localhost"
	port             = "8080"
	timeoutStr       = "TIMEOUT_SECONDS"
	maxRetryAttempts = 5
)

func main() {
	// Set up configurable timeout
	timeout, err := getConfigurableTimeout()
	if err != nil {
		log.Fatalf("Invalid configurable timeout: %v", err)
	}
	dialer := &net.Dialer{Timeout: timeout}

	for i := 0; i < maxRetryAttempts; i++ {
		conn, err := dialer.Dial("tcp", host+":"+port)
		if err != nil {
			handleError("Connection failed", err)
			continue
		}
		defer conn.Close()
		if err := sendRequest(conn); err != nil {
			handleError("Sending request failed", err)
			continue
		}
		if err := readResponse(conn); err != nil {
			handleError("Reading response failed", err)
			continue
		}
		fmt.Println("Success!")
		break
	}
}
func getConfigurableTimeout() (time.Duration, error) {
	timeoutStr := os.Getenv(timeoutStr)
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return 0, fmt.Errorf("invalid timeout environment variable: %w", err)
	}
	if timeout <= 0 {
		return 0, fmt.Errorf("timeout must be a positive integer")
	}
	return time.Duration(timeout) * time.Second, nil
}
func sendRequest(conn net.Conn) error {
	request := "Hello, World!"
	_, err := conn.Write([]byte(request))
	return err
}
func readResponse(conn net.Conn) error {
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return err
	}
	fmt.Println(string(response[:n]))
	return nil
}

func handleError(msg string, err error) {
	if nerr, ok := err.(*net.OpError); ok && nerr.Timeout() {
		fmt.Printf("%s - Timeout, retrying...\n", msg)
		backoffTime := time.Duration(len(msg)) * time.Second // Implement exponential backoff here if needed
		time.Sleep(backoffTime)
		return
	}
	log.Fatalf("%s: %v", msg, err)
}
