package _94209

import (
	"fmt"
	"testing"
)

func singleReturn(a, b int) int {
	return a + b
}

func multipleReturn(a, b int) (int, int) {
	sum := a + b
	return sum, sum * 2
}

func BenchmarkSingleReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		singleReturn(10, 20)
	}
}

func BenchmarkMultipleReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = multipleReturn(10, 20)
	}
}

func main() {
	fmt.Println("Single return value:", singleReturn(10, 20))
	sum, doubleSum := multipleReturn(10, 20)
	fmt.Println("Multiple return values:", sum, doubleSum)
}
