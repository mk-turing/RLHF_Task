package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var availableVariants = []string{"a", "b", "c"}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleRequest)
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	queryParams, _ := url.ParseQuery(r.URL.Query().Encode())

	// Determine the variant from cookie first
	variantCookie, err := r.Cookie("variant")
	if err == nil && variantCookie.Value != "" {
		// Use variant from cookie if available
		variant := variantCookie.Value
		serveVariant(w, r, variant)
		return
	}

	// Determine the variant from query parameter
	variant := queryParams.Get("variant")
	if variant != "" {
		if strings.Contains(strings.Join(availableVariants, ","), variant) {
			// Valid variant from query parameter
			serveVariant(w, r, variant)
			return
		} else {
			http.Error(w, "Invalid variant parameter", http.StatusBadRequest)
			return
		}
	}

	// Assign a random variant if no cookie or query parameter is specified
	variant = availableVariants[rand.Intn(len(availableVariants))]

	// Create a cookie for the variant
	cookie := http.Cookie{
		Name:     "variant",
		Value:    variant,
		Path:     "/",
		MaxAge:   3600 * 24 * 30, // 30 days
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	// Serve the variant
	serveVariant(w, r, variant)
}

func serveVariant(w http.ResponseWriter, r *http.Request, variant string) {
	switch variant {
	case "a":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant A</h1><p>This is the content for variant A.</p></body></html>")
	case "b":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant B</h1><p>This is the content for variant B.</p></body></html>")
	case "c":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant C</h1><p>This is the content for variant C.</p></body></html>")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1><p>Invalid variant parameter.</p></body></html>")
	}
}
