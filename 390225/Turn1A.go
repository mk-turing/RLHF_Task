package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

func serviceA(w http.ResponseWriter, r *http.Request) {
	// Simulate service A processing
	// This could be replaced with actual service calls
	response := "Response from Service A"
	fmt.Fprintln(w, response)
}

func serviceB(w http.ResponseWriter, r *http.Request) {
	// Simulate service B processing
	response := "Response from Service B"
	fmt.Fprintln(w, response)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	results := make([]string, 2)

	// Invoke Service A
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := captureResponse(serviceA) // Capture the response from service A
		results[0] = res
	}()

	// Invoke Service B
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := captureResponse(serviceB) // Capture the response from service B
		results[1] = res
	}()

	wg.Wait() // Wait for all goroutines to finish

	// Aggregate results
	for _, result := range results {
		fmt.Fprintln(w, result)
	}
}

// Helper function to capture the response from a service function
func captureResponse(handler http.HandlerFunc) string {
	rec := httptest.NewRecorder() // Mocking HTTP response
	req := httptest.NewRequest("GET", "/", nil)
	handler(rec, req) // Call the service
	resp := rec.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	http.HandleFunc("/gateway", handleRequest) // API Gateway endpoint
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
