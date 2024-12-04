//package main
//
//import (
//	"context"
//	"log"
//	"net"
//	"sync"
//	"time"
//)
//
//const (
//	serverAddr = "localhost:8080"
//	poolSize   = 2 // Adjust the pool size as needed based on your requirements
//)
//
//// connectionPool manages a pool of TCP connections
//type connectionPool struct {
//	mutex       sync.Mutex
//	connections []net.Conn
//}
//
//func newConnectionPool() *connectionPool {
//	return &connectionPool{
//		connections: make([]net.Conn, 0, poolSize),
//	}
//}
//
//func (cp *connectionPool) getConnection(ctx context.Context) (net.Conn, error) {
//	cp.mutex.Lock()
//	defer cp.mutex.Unlock()
//
//	for _, conn := range cp.connections {
//		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
//			log.Println("Setting read deadline failed:", err)
//			continue
//		}
//
//		_, err := conn.Read(make([]byte, 1))
//		if err == nil {
//			return conn, nil
//		}
//		log.Println("Connection inactive:", err)
//	}
//
//	// If no active connection is found, create a new one
//	conn, err := net.Dial("tcp", serverAddr)
//	if err != nil {
//		return nil, err
//	}
//	cp.connections = append(cp.connections, conn)
//	return conn, nil
//}
//
//func (cp *connectionPool) releaseConnection(conn net.Conn) {
//	cp.mutex.Lock()
//	defer cp.mutex.Unlock()
//	cp.connections = append(cp.connections, conn)
//}
//
//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	pool := newConnectionPool()
//
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			default:
//				sendData(ctx, pool)
//				time.Sleep(time.Second * 2) // Simulate periodic data sending
//			}
//		}
//	}()
//
//	select {
//	case <-ctx.Done():
//	}
//}
//
//func sendData(ctx context.Context, pool *connectionPool) {
//	conn, err := pool.getConnection(ctx)
//	if err != nil {
//		log.Println("Error getting connection:", err)
//		return
//	}
//	defer pool.releaseConnection(conn)
//
//	message := []byte("Hello, Server!")
//	_, err = conn.Write(message)
//	if err != nil {
//		log.Println("Error writing data:", err)
//		return
//	}
//
//	log.Println("Data sent successfully.")
//}
