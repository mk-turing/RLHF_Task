//package main
//
//import (
//	"context"
//	"fmt"
//	"log"
//	"net"
//	"net/http"
//	"os"
//	"os/signal"
//	"sync"
//	"syscall"
//	"time"
//
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus_http"
//)
//
//const (
//	serverAddress = ":8080"
//	connectionTimeout  = 5 * time.Second
//	readTimeout        = 10 * time.Second
//	writeTimeout       = 10 * time.Second
//	keepAliveInterval = 3 * time.Minute
//)
//
//// Metrics variables
//var (
//	requestsReceived = prometheus.NewCounterVec(
//		prometheus.CounterOpts{
//			Name: "server_requests_received",
//			Help: "Total number of requests received by the server.",
//		},
//		[]string{"method"},
//	)
//	requestProcessingTime = prometheus.NewSummaryVec(
//		prometheus.SummaryOpts{
//			Name: "server_request_processing_seconds",
//			Help: "Processing time of requests by the server in seconds.",
//		},
//		[]string{"method"},
//	)
//)
//
//func init() {
//	// Register metrics
//	prometheus.MustRegister(requestsReceived)
//	prometheus.MustRegister(requestProcessingTime)
//
//	// Start a metrics HTTP server
//	go func() {
//		log.Fatal(http.ListenAndServe(":9100", prometheus_http.Handler()))
//	}()
//}
//
//var (
//	wg sync.WaitGroup
//)
//
//func main() {
//	log.Printf("Server starting on %s", serverAddress)
//
//	go handleSignals()
//
//	ln, err := net.Listen("tcp", serverAddress)
//	if err != nil {
//		log.Fatalf("Error listening: %v", err)
//	}
//
//	defer ln.Close()
//
//	for {
//		conn, err := ln.Accept()
//		if err != nil {
//			if os.IsInterrupt(err) {
//				log.Printf("Server shutting down...")
//				return
//			}
//			log.Printf("Error accepting connection: %v", err)
//			continue
//		}
//
//		wg.Add(1)
//		go handleConnection(conn)
//	}
//}
//
//func handleSignals() {
//	signals := make(chan os.Signal)
//	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
//	for _ = range signals {
//		log.Printf("Received shutdown signal, closing connections...")
//		// Close existing connections
//		wg.Wait()
//		log.Printf("Server shut down")
//		os.Exit(0)
//	}
//}
//
//func handleConnection(conn net.Conn) {
//	defer wg.Done()
//
//	defer conn.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)