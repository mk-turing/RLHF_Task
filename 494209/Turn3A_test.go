package _94209

import (
	"errors"
	"testing"
)

func sumSingle(numbers *[]int) int {
	total := 0
	for _, num := range *numbers {
		total += num
	}
	return total
}

func sumAndAverage(numbers *[]int) (int, float64, error) {
	total := 0
	count := 0
	for _, num := range *numbers {
		total += num
		count++
	}
	if count == 0 {
		return 0, 0, errors.New("divide by zero")
	}
	return total, float64(total) / float64(count), nil
}

func sumAndAverageErrorFirst(numbers *[]int) (error, int, float64) {
	total := 0
	count := 0
	for _, num := range *numbers {
		total += num
		count++
	}
	if count == 0 {
		return errors.New("divide by zero"), 0, 0
	}
	return nil, total, float64(total) / float64(count)
}

type Result struct {
	Sum     int
	Average float64
}

func sumAndAverageStruct(numbers *[]int) (*Result, error) {
	total := 0
	count := 0
	for _, num := range *numbers {
		total += num
		count++
	}
	if count == 0 {
		return nil, errors.New("divide by zero")
	}
	result := &Result{Sum: total, Average: float64(total) / float64(count)}
	return result, nil
}

func BenchmarkSumSingle(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}
	numbersPtr := &numbers

	for n := 0; n < b.N; n++ {
		sumSingle(numbersPtr)
	}
}

func BenchmarkSumAndAverage(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}
	numbersPtr := &numbers

	for n := 0; n < b.N; n++ {
		_, _, _ = sumAndAverage(numbersPtr)
	}
}

func BenchmarkSumAndAverageErrorFirst(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}
	numbersPtr := &numbers

	for n := 0; n < b.N; n++ {
		_, _, _ = sumAndAverageErrorFirst(numbersPtr)
	}
}

func BenchmarkSumAndAverageStruct(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}
	numbersPtr := &numbers

	for n := 0; n < b.N; n++ {
		_, _ = sumAndAverageStruct(numbersPtr)
	}
}
