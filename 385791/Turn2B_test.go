package main

import (
	"math/rand"
	"testing"
)

func BenchmarkReverseStringRandom(b *testing.B) {
	// Benchmarking a randomly generated string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Seed(int64(i))
		b.SetBytes(int64(len(randomString(rand.Intn(200) + 1))))
		reverseString(randomString(rand.Intn(200) + 1))
	}
}

func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s := []byte(letters)
	b := make([]byte, length)
	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}
	return string(b)
}
