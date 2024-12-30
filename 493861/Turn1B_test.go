package main

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// using time.Now()
func BenchmarkTimeNow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		time.Now()
	}
}

// using time.Since()
func BenchmarkTimeSince(b *testing.B) {
	start := time.Now()
	for n := 0; n < b.N; n++ {
		// some computation here
	}
	time.Since(start)
}

// using runtime.Now()
func BenchmarkRuntimeNow(b *testing.B) {
	for n := 0; n < b.N; n++ {
		time.Now()
	}
}

// using runtime.NumGoroutine()
func BenchmarkRuntimeNumGoroutine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runtime.NumGoroutine()
	}
}

// using sync.WaitGroup
func BenchmarkSyncWaitGroup(b *testing.B) {
	wg := &sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
		wg.Wait()
	}
}

// using time.Duration
func BenchmarkTimeDuration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var d time.Duration
		d += 1 * time.Second
		d += 1 * time.Nanosecond
		d += 1 * time.Microsecond
		d += 1 * time.Millisecond
		d += 1 * time.Minute
		d += 1 * time.Hour
		d += 1 * time.Second
		d += 1 * time.Nanosecond
		d += 1 * time.Microsecond
		d += 1 * time.Millisecond
		d += 1 * time.Minute
		d += 1 * time.Hour
		d += 1 * time.Second
		d += 1 * time.Nanosecond
		d += 1 * time.Microsecond
		d += 1 * time.Millisecond
		d += 1 * time.Minute
		d += 1 * time.Hour
		d += 1 * time.Second
		d += 1 * time.Nanosecond
		d += 1 * time.Microsecond
		d += 1 * time.Millisecond
		d += 1 * time.Minute
		d += 1 * time.Hour
	}
}

// using time.Ticker
func BenchmarkTimeTicker(b *testing.B) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for n := 0; n < b.N; n++ {
		<-ticker.C
	}
}

// using time.After()
func BenchmarkTimeAfter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		<-time.After(1 * time.Second)
	}
}

// using time.Sleep()
func BenchmarkTimeSleep(b *testing.B) {
	for n := 0; n < b.N; n++ {
		time.Sleep(1 * time.Second)
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	testing.Benchmark(BenchmarkTimeNow)
	testing.Benchmark(BenchmarkTimeSince)
	testing.Benchmark(BenchmarkRuntimeNow)
	testing.Benchmark(BenchmarkRuntimeNumGoroutine)
	testing.Benchmark(BenchmarkSyncWaitGroup)
	testing.Benchmark(BenchmarkTimeDuration)
	testing.Benchmark(BenchmarkTimeTicker)
	testing.Benchmark(BenchmarkTimeAfter)
	testing.Benchmark(BenchmarkTimeSleep)
}
