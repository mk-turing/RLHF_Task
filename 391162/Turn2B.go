package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	url                = "http://localhost:8080"
	totalRequests      = 10000            // Total number of requests
	concurrentRequests = 100              // Number of concurrent requests
	duration           = 30 * time.Second // Test duration
)

var (
	startTime     = time.Time{}
	responseTimes []time.Duration
	errorCount    = int32(0)
)

func benchmark() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	var queue = make(chan struct{}, concurrentRequests)

	startTime = time.Now()

	for i := 0; i < totalRequests; i++ {
		queue <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-queue
			}()
			start := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				return
			}
			_, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			responseTimes = append(responseTimes, time.Since(start))
		}()
	}

	wg.Wait()

	totalTime := time.Since(startTime)
	// Sum the response times manually
	var totalResponseTime time.Duration
	for _, rt := range responseTimes {
		totalResponseTime += rt
	}
	avgResponseTime := float64(totalResponseTime) / float64(len(responseTimes))
	throughput := float64(len(responseTimes)) / totalTime.Seconds()

	errorRate := float64(atomic.LoadInt32(&errorCount)) / float64(totalRequests) * 100

	fmt.Printf("Total Time: %.2f seconds\n", totalTime.Seconds())
	fmt.Printf("Average Response Time: %.2f ms\n", avgResponseTime)
	fmt.Printf("Throughput: %.2f requests/sec\n", throughput)
	fmt.Printf("Error Rate: %.2f%%\n", errorRate)

	// Sort response times
	sort.Slice(responseTimes, func(i, j int) bool {
		return responseTimes[i] < responseTimes[j]
	})

	// Calculate response time percentiles
	fmt.Printf("P50 Response Time: %.2f ms\n", float64(responseTimes[len(responseTimes)/2])/float64(time.Millisecond))
	fmt.Printf("P90 Response Time: %.2f ms\n", float64(responseTimes[9*len(responseTimes)/10])/float64(time.Millisecond))
	fmt.Printf("P99 Response Time: %.2f ms\n", float64(responseTimes[99*len(responseTimes)/100])/float64(time.Millisecond))
}

func main() {
	benchmark()
}
