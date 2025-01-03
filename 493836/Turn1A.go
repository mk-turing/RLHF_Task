package main

import (
	"fmt"
	"net/http"
	"sync"
)

var tenantChannels map[string]chan string = make(map[string]chan string)
var wg sync.WaitGroup

func handler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	event := "New request received"
	if channel, exists := tenantChannels[tenantID]; exists {
		channel <- event
	} else {
		// Create a new channel for this tenant
		tenantChannels[tenantID] = make(chan string, 100)
		go processEvents(tenantID)
		tenantChannels[tenantID] <- event
	}

	fmt.Fprintf(w, "Hello, tenant %s! Your event: %s", tenantID, event)
}

func processEvents(tenantID string) {
	for event := range tenantChannels[tenantID] {
		fmt.Printf("Processing event for tenant %s: %s\n", tenantID, event)
	}
	delete(tenantChannels, tenantID)
	wg.Done()
}

func main() {
	http.HandleFunc("/", handler)
	wg.Add(1)
	go func() {
		wg.Wait()
		fmt.Println("Server shutting down")
	}()
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
