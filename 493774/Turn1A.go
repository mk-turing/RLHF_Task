package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	queryParams, _ := url.ParseQuery(r.URL.Query().Encode())

	// Determine the variant
	variant := queryParams.Get("variant")
	if variant == "" {
		variant = "a" // Default variant
	}

	// Serve different content based on variant
	switch variant {
	case "a":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant A</h1><p>This is the content for variant A.</p></body></html>")
	case "b":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant B</h1><p>This is the content for variant B.</p></body></html>")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1><p>Invalid variant parameter.</p></body></html>")
	}
}
