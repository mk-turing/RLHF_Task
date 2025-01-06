package _94209

import "testing"

func sumSingle(numbers []int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

func sumAndAverage(numbers []int) (int, float64) {
	total := 0
	count := 0
	for _, num := range numbers {
		total += num
		count++
	}
	if count == 0 {
		return 0, 0
	}
	return total, float64(total) / float64(count)
}

func BenchmarkSumSingle(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		sumSingle(numbers)
	}
}

func BenchmarkSumAndAverage(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		sumAndAverage(numbers)
	}
}
