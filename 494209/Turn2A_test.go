package _94209

import "testing"

type CalcError error

func (e CalcError) Error() string {
	return e.Error()
}

func sumSingle(numbers []int) (int, error) {
	total := 0
	for _, num := range numbers {
		if num > 1000 {
			return 0, CalcError("Number exceeds threshold")
		}
		total += num
	}
	return total, nil
}

func sumAndAverage(numbers []int) (int, float64, error) {
	total := 0
	count := 0
	for _, num := range numbers {
		if num > 1000 {
			return 0, 0, CalcError("Number exceeds threshold")
		}
		total += num
		count++
	}
	if count == 0 {
		return 0, 0, CalcError("Division by zero")
	}
	return total, float64(total) / float64(count), nil
}

func BenchmarkSumSingleNoError(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i % 1000
	}

	for n := 0; n < b.N; n++ {
		sumSingle(numbers)
	}
}

func BenchmarkSumAndAverageNoError(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i % 1000
	}

	for n := 0; n < b.N; n++ {
		sumAndAverage(numbers)
	}
}

func BenchmarkSumSingleWithError(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		sum, err := sumSingle(numbers)
		if err != nil {
			b.Fatal(err) // Fail on error to ensure we are not benchmarking error handling itself
		}
	}
}

func BenchmarkSumAndAverageWithError(b *testing.B) {
	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = i
	}

	for n := 0; n < b.N; n++ {
		sum, avg, err := sumAndAverage(numbers)
		if err != nil {
			b.Fatal(err)
		}
	}
}
