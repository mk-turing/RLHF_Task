package main

import (
	"testing"
)

func someFunction(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result += i * i
	}
	return result
}

func BenchmarkSomeFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		someFunction(10000) // Input size can be adjusted as needed.
	}
}
