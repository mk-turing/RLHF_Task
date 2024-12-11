package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	url                = "http://localhost:8080"
	totalRequests      = 1000 // Total number of requests
	concurrentRequests = 100  // Number of concurrent requests
)

func benchmark() {
	var wg sync.WaitGroup
	responseTimes := make([]time.Duration, totalRequests)

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			startTime := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			_, _ = ioutil.ReadAll(resp.Body) // Read and discard the response body
			resp.Body.Close()
			responseTimes[i] = time.Since(startTime)
		}(i)

		if (i+1)%concurrentRequests == 0 {
			wg.Wait() // Wait for the current batch to finish before starting the next
		}
	}

	wg.Wait() // Wait for any remaining requests to finish

	elapsed := time.Since(start)
	totalTime := float64(elapsed.Milliseconds())
	avgResponseTime := totalTime / float64(totalRequests)
	throughput := float64(totalRequests) / totalTime * 1000 // Requests per second

	fmt.Printf("Total Time: %.2f ms\n", totalTime)
	fmt.Printf("Average Response Time: %.2f ms\n", avgResponseTime)
	fmt.Printf("Throughput: %.2f requests/sec\n", throughput)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(50 * time.Millisecond) // Simulate work
	fmt.Fprintf(w, "Hello, World!")
}

func benchmarkHelper() {
	time.Sleep(5 * time.Second)
	benchmark()
}

func main() {
	go benchmarkHelper()
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8080", nil)
}
