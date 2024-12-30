package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"testing"
	"time"
)

// Sorting function
func sortArray(arr []int) {
	sort.Ints(arr)
}

// Searching function
func searchElement(arr []int, key int) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

// Mathematical computation function (Factorial)
func factorial(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

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

// using monotonic clock (replacing syscall.Gettimeofday)
func BenchmarkMonotonicClock(b *testing.B) {
	for n := 0; n < b.N; n++ {
		start := time.Now()
		_ = start.UnixNano() // Use UnixNano as a high-precision timestamp
	}
}

// Benchmark for sorting
func BenchmarkSortArray(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = rand.Intn(size)
		}
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				sortArray(arr)
			}
		})
	}
}

// Benchmark for searching
func BenchmarkSearchElement(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = rand.Intn(size)
		}
		key := arr[rand.Intn(size)]
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				searchElement(arr, key)
			}
		})
	}
}

// Benchmark for mathematical computation (Factorial)
func BenchmarkFactorial(b *testing.B) {
	for _, n := range []int{10, 20, 30, 40} {
		b.Run(fmt.Sprintf("n-%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				factorial(n)
			}
		})
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	testing.Benchmark(BenchmarkTimeNow)
	testing.Benchmark(BenchmarkTimeSince)
	testing.Benchmark(BenchmarkMonotonicClock)
	testing.Benchmark(BenchmarkSortArray)
	testing.Benchmark(BenchmarkSearchElement)
	testing.Benchmark(BenchmarkFactorial)
}
