package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"
)

// Simulated service responses
var cache = map[string]string{
	"serviceA": "Cached response from Service A",
	"serviceB": "Cached response from Service B",
}

func serviceA(w http.ResponseWriter, r *http.Request) {
	if response, ok := cache["serviceA"]; ok {
		sendResponse(w, response)
		return
	}

	// Simulate service A processing
	time.Sleep(100 * time.Millisecond)
	response := "Response from Service A"
	cache["serviceA"] = response // Cache the response
	sendResponse(w, response)
}

func serviceB(w http.ResponseWriter, r *http.Request) {
	if response, ok := cache["serviceB"]; ok {
		sendResponse(w, response)
		return
	}

	// Simulate service B processing
	time.Sleep(100 * time.Millisecond)
	response := "Response from Service B"
	cache["serviceB"] = response // Cache the response
	sendResponse(w, response)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	results := make([]string, 2)

	// Invoke Service A
	wg.Add(1)
	go func() {
		defer wg.Done()
		results[0] = captureResponse(serviceA)
	}()

	// Invoke Service B
	wg.Add(1)
	go func() {
		defer wg.Done()
		results[1] = captureResponse(serviceB)
	}()

	wg.Wait()

	// Aggregate results
	response, _ := json.Marshal(results)
	sendResponse(w, string(response))
}

func sendResponse(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	fmt.Fprintf(gz, response)
}

func captureResponse(handler http.HandlerFunc) string {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler(rec, req)
	resp := rec.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	var requestCount atomic.Int64

	http.HandleFunc("/gateway", func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		handleRequest(w, r)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	go func() {
		for range time.Tick(time.Second) {
			count := requestCount.Load()
			log.Printf("Handled %d requests in the last second\n", count)
			requestCount.Store(0)
		}
	}()
}
