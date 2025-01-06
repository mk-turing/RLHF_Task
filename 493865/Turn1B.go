package main

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkFunc represents the function to be benchmarked.
type BenchmarkFunc func()

// Benchmark performs the benchmarking and prints the result.
func Benchmark(b *testing.B, name string, f BenchmarkFunc) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
	b.StopTimer()
	fmt.Printf("%s: %s\n", name, b.ReportAllocs())
}

func ExampleBenchmark() {
	Benchmark(testing.Benchmark, "Sleep100us", func() {
		time.Sleep(100 * time.Microsecond)
	})
	// Output: Sleep100us: 33.3 ns/op    0 B/op    0 allocs/op
}

func main() {
	ExampleBenchmark()
}
