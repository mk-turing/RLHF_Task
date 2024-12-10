package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func callServiceWithRetries(w http.ResponseWriter, r *http.Request, serviceURL string, retries int) {
	// Maximum backoff time of 10 seconds
	maxBackoff := 10 * time.Second

	for i := 0; i < retries; i++ {
		resp, err := http.Get(serviceURL)
		if err != nil {
			log.Printf("Attempt %d: Error calling service: %v. Retrying...", i+1, err)
			// Calculate the exponential backoff time, capped at maxBackoff
			backoffDuration := time.Duration(1<<i) * time.Second // Exponential backoff (2^i seconds)
			if backoffDuration > maxBackoff {
				backoffDuration = maxBackoff
			}
			log.Printf("Retrying in %v...", backoffDuration)
			time.Sleep(backoffDuration)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Attempt %d: Service returned bad status: %v. Retrying...", i+1, resp.Status)
			// Exponential backoff with cap on maxBackoff
			backoffDuration := time.Duration(1<<i) * time.Second
			if backoffDuration > maxBackoff {
				backoffDuration = maxBackoff
			}
			log.Printf("Retrying in %v...", backoffDuration)
			time.Sleep(backoffDuration)
			continue
		}

		// Successful response
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintln(w, string(body))
		return
	}

	// If we reach here, all retries failed
	fmt.Fprintln(w, "Failed to reach service after all attempts.")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	urlA := "http://localhost:8081/serviceA"
	urlB := "http://localhost:8082/serviceB"

	callServiceWithRetries(w, r, urlA, 3)
	callServiceWithRetries(w, r, urlB, 3)
}

func main() {
	http.HandleFunc("/gateway", handleRequest)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
