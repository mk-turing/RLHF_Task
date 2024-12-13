package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func fetch(url string, wg *sync.WaitGroup, latency chan<- time.Duration) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	end := time.Now()
	latency <- end.Sub(start)
}

func main() {
	var wg sync.WaitGroup
	latency := make(chan time.Duration)

	urls := []string{
		"https://example.com",
		"https://example.com",
		"https://example.com",
		"https://example.com",
		"https://example.com",
	}

	for _, url := range urls {
		wg.Add(1)
		go fetch(url, &wg, latency)
	}

	wg.Wait()

	close(latency)
	fmt.Println("Latency statistics:")
	for l := range latency {
		fmt.Printf("%v\n", l)
	}
}
