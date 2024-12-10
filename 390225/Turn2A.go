package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	URL string
}

func (s *Service) Call() (string, error) {
	resp, err := http.Get(s.URL) // Reusable Get request
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

// API Gateway Handler
func apiGatewayHandler(w http.ResponseWriter, r *http.Request) {
	services := []*Service{
		{URL: "http://localhost:8081/serviceA"},
		{URL: "http://localhost:8082/serviceB"},
	}

	var wg sync.WaitGroup
	responses := make([]string, len(services))
	errs := make([]error, len(services))

	for i, service := range services {
		wg.Add(1)
		go func(i int, service *Service) {
			defer wg.Done()
			response, err := service.Call()
			responses[i] = response
			errs[i] = err
		}(i, service)
	}

	wg.Wait() // Wait for all services to respond

	for i, err := range errs {
		if err == nil {
			fmt.Fprintf(w, "Response from service %d: %s\n", i+1, responses[i])
		} else {
			fmt.Fprintf(w, "Error from service %d: %v\n", i+1, err)
		}
	}
}

func main() {
	http.HandleFunc("/gateway", apiGatewayHandler) // API Gateway endpoint
	go func() {
		log.Fatal(http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock response for service A
			time.Sleep(50 * time.Millisecond) // Simulate delay
			fmt.Fprintln(w, "This is Service A")
		})))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8082", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock response for service B
			time.Sleep(30 * time.Millisecond) // Simulate delay
			fmt.Fprintln(w, "This is Service B")
		})))
	}()

	log.Println("Starting API Gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
