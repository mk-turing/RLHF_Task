//package main
//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	"log"
//	"net"
//	"sync"
//	"time"
//)
//
//const (
//	serverAddress     = "localhost:8080"
//	maxConnections    = 100
//	connectionTimeout = 5 * time.Second
//	readTimeout       = 10 * time.Second
//	writeTimeout      = 10 * time.Second
//	keepAliveInterval = 3 * time.Minute
//)
//
//var (
//	// Use a global sync.WaitGroup to manage goroutines
//	wg sync.WaitGroup
//)
//
//func main() {
//	// Create a connection pool using a sync.Pool
//	connectionPool := &sync.Pool{
//		New: func() interface{} {
//			conn, err := net.DialTimeout("tcp", serverAddress, connectionTimeout)
//			if err != nil {
//				log.Fatalf("Error dialing: %v", err)
//			}
//
//			// Set keep-alive on the connection
//			tcpConn, ok := conn.(*net.TCPConn)
//			if ok {
//				if err := tcpConn.SetKeepAlive(true); err != nil {
//					log.Printf("Error setting keep-alive: %v", err)
//				}
//				if err := tcpConn.SetKeepAlivePeriod(keepAliveInterval); err != nil {
//					log.Printf("Error setting keep-alive period: %v", err)
//				}
//			}
//
//			// Set read and write timeouts on the connection
//			//conn.SetReadTimeout(readTimeout)
//			//conn.SetWriteTimeout(writeTimeout)
//
//			return conn
//		},
//	}
//
//	defer func() {
//		// Close all connections in the pool when the application exits
//		for {
//			if conn, ok := connectionPool.Get().(net.Conn); ok {
//				conn.Close()
//				connectionPool.Put(nil)
//			} else {
//				break
//			}
//		}
//	}()
//
//	// Start multiple goroutines to handle concurrent connections
//	for i := 0; i < maxConnections; i++ {
//		wg.Add(1)
//		go handleConnection(connectionPool)
//	}
//
//	wg.Wait()
//}
//
//func handleConnection(connectionPool *sync.Pool) {
//	defer wg.Done()
//
//	for {
//		// Get a connection from the pool
//		conn, ok := connectionPool.Get().(net.Conn)
//		if !ok {
//			// No connections available in the pool, exit goroutine
//			return
//		}
//
//		defer conn.Close()
//
//		_, cancel := context.WithTimeout(context.Background(), time.Minute)
//		defer cancel()
//
//		// Perform read and write operations using buffering
//		br := bufio.NewReader(conn)
//		bw := bufio.NewWriter(conn)
//
//		// Read data from the server
//		data, err := br.ReadBytes('\n')
//		if err != nil {
//			if err == context.DeadlineExceeded {
//				log.Printf("Read timeout exceeded")
//			} else {
//				log.Printf("Error reading: %v", err)
//			}
//			continue
//		}
//
//		// Process the received data
//		processedData := processData(string(data))
//
//		// Write data back to the server
//		if _, err := bw.WriteString(processedData + "\n"); err != nil {
//			if err == context.DeadlineExceeded {
//				log.Printf("Write timeout exceeded")
//			} else {
//				log.Printf("Error writing: %v", err)
//			}
//			continue
//		}
//
//		if err := bw.Flush(); err != nil {
//			log.Printf("Error flushing buffer: %v", err)
//			continue
//		}
//	}
//}
//
//func processData(input string) string {
//	// Simulate data processing
//	time.Sleep(1 * time.Second)
//	return fmt.Sprintf("Processed: %s", input)
//}
