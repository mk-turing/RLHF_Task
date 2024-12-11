package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func TestHelloHandler(t *testing.T) {
	// Create a new router
	r := http.NewServeMux()
	r.HandleFunc("/", HelloHandler)

	// Start a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Benchmark the handler
	t.Run("BenchmarkHelloHandler", func(t *testing.T) {
		client := &http.Client{}
		url := ts.URL + "/"

		// Measure response time
		t.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := client.Get(url)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
