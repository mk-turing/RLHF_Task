package main

import (
	"sort"
	"testing"
	"time"
)

// Using time.Now()
func benchmarkTimeNow(b *testing.B, f func()) {
	for n := 0; n < b.N; n++ {
		start := time.Now()
		f()
		time.Since(start)
	}
}

// Using runtime.Nanoseconds()
func benchmarkRuntimeNanoseconds(b *testing.B, f func()) {
	for n := 0; n < b.N; n++ {
		f()
	}
}

// Using monotonicClock.Now()
func benchmarkMonotonicClock(b *testing.B, f func()) {
	for n := 0; n < b.N; n++ {
		f()
	}
}

// Short-running function for sorting a slice of integers
func sortShort(a []int) {
	sort.Ints(a)
}

// Long-running function for sorting a larger slice of integers
func sortLong(a []int) {
	for i := range a {
		for j := 0; j < len(a)-i-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
			}
		}
	}
}

// Short-running function for performing a mathematical computation
func mathShort() float64 {
	sum := 0.0
	for i := 0; i < 1000; i++ {
		sum += 1.0 / float64(i+1)
	}
	return sum
}

// Long-running function for performing a mathematical computation
func mathLong() float64 {
	sum := 0.0
	for i := 0; i < 1000000; i++ {
		sum += 1.0 / float64(i+1)
	}
	return sum
}

// Short-running function for searching for a value in a slice
func searchShort(a []int, target int) bool {
	for _, val := range a {
		if val == target {
			return true
		}
	}
	return false
}

// Long-running function for searching for a value in a larger slice
func searchLong(a []int, target int) bool {
	for _, val := range a {
		if val == target {
			return true
		}
	}
	return false
}

func main() {
	shortArray := []int{4, 1, 3, 2}
	longArray := make([]int, 10000)
	for i := range longArray {
		longArray[i] = i
	}

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { sortShort(shortArray) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { sortShort(shortArray) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { sortShort(shortArray) })
	})

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { sortLong(longArray) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { sortLong(longArray) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { sortLong(longArray) })
	})

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { mathShort() })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { mathShort() })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { mathShort() })
	})

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { mathLong() })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { mathLong() })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { mathLong() })
	})

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { searchShort(shortArray, 1) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { searchShort(shortArray, 1) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { searchShort(shortArray, 1) })
	})

	testing.Benchmark(func(b *testing.B) {
		benchmarkTimeNow(b, func() { searchLong(longArray, 1) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkRuntimeNanoseconds(b, func() { searchLong(longArray, 1) })
	})
	testing.Benchmark(func(b *testing.B) {
		benchmarkMonotonicClock(b, func() { searchLong(longArray, 1) })
	})
}
