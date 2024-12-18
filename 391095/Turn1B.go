package main
import (
	"fmt"
	"log"
	"net"
	"time"
)

// Config struct for configurable timeout strategy
type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewClientWithConfig creates a new TCP client with configurable timeout
func NewClientWithConfig(addr string, config Config) (*net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, config.ReadTimeout)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now().Add(config.WriteTimeout))
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// NewClientWithHardCodedTimeout creates a new TCP client with hard-coded timeout
func NewClientWithHardCodedTimeout(addr string) (*net.Conn, error) {
	hardCodedReadTimeout := 5 * time.Second
	hardCodedWriteTimeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", addr, hardCodedReadTimeout)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now().Add(hardCodedWriteTimeout))
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func main() {
	serverAddr := "localhost:8080"

	// Configurability
	config := Config{
		ReadTimeout:  10 * time.Second, // Adjustable read timeout
		WriteTimeout: 5 * time.Second,  // Adjustable write timeout
	}
	conn, err := NewClientWithConfig(serverAddr, config)
	if err != nil {
		log.Fatal("Error connecting to server: ", err)