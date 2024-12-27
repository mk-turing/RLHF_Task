package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	variantCookieName = "abctest"
	defaultVariant    = "a"
)

var (
	variants = []string{"a", "b", "c"}
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	variant := getVariantFromCookie(r)
	if variant == "" {
		variant = chooseRandomVariant()
		setVariantCookie(w, variant)
	}
	serveVariantContent(w, variant)
}

func getVariantFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(variantCookieName)
	if err == nil {
		return cookie.Value
	}
	return ""
}

func setVariantCookie(w http.ResponseWriter, variant string) {
	cookie := http.Cookie{
		Name:     variantCookieName,
		Value:    variant,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7), // Set cookie expiration for 7 days
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func chooseRandomVariant() string {
	// Select a random variant from the available variants
	return variants[rand.Intn(len(variants))]
}

func serveVariantContent(w http.ResponseWriter, variant string) {
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
