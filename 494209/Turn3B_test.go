package _94209

import (
	"errors"
	"testing"
)

// Refactored ResultStruct to hold sum and average
type ResultStruct struct {
	Sum     int
	Average float64
}

func sumSingle(numbers []int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

func sumAndAveragePtr(numbers []int) (*ResultStruct, error) {
	total := 0
	count := 0
	for _, num := range numbers {
		total += num
		count++
	}

	if count == 0 {
		return nil, errors.New("divide by zero")
	}

	result := &ResultStruct{Sum: total, Average: float64(total) / float64(count)}
	return result, nil
}

func sumAndAverageStruct(numbers []int) (ResultStruct, error) {
	total := 0
	count := 0
	for _, num := range numbers {
		total += num
		count++
	}

	if count == 0 {
		return ResultStruct{}, errors.New("divide by zero")
	}

	result := &ResultStruct{Sum: total, Average: float64(total) / float64(count)}
	return *result, nil
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

func BenchmarkSumAndAveragePtr(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		_, _ = sumAndAveragePtr(numbers)
	}
}

func BenchmarkSumAndAverageStruct(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		_, _ = sumAndAverageStruct(numbers)
	}
}
