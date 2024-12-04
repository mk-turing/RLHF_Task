package main

import "testing"

func reverseString(input string) string {
	reversed := []byte(input)
	for i := 0; i < len(reversed)/2; i++ {
		reversed[i], reversed[len(reversed)-i-1] = reversed[len(reversed)-i-1], reversed[i]
	}
	return string(reversed)
}

func BenchmarkReverseString(b *testing.B) {
	input := "Hello World"
	for i := 0; i < b.N; i++ {
		reverseString(input)
	}
}
