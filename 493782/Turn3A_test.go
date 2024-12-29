package main

import (
	"benchmark"
	"errors"
	"sync"
	"time"
)

// processDataWithConcurrency processes data using multiple goroutines and calls callbacks accordingly.
func processDataWithConcurrency(data [][]byte, callbacks []func(error)) {
	wg := sync.WaitGroup{}
	wg.Add(len(data))

	for i, d := range data {
		go func(idx int, data []byte, callback func(error)) {
			defer wg.Done()

			// Simulated processing
			if data == nil {
				callback(errors.New("data is nil"))
				return
			}

			// Simulate error after a random delay to simulate concurrent execution
			time.Sleep(time.Duration(idx*100) * time.Millisecond)
			if len(data) == 0 {
				callback(errors.New("data is empty"))
				return
			}

			callback(nil)
		}(i, d, callbacks[i])
	}

	wg.Wait()
}

// MockCallback saves the error message in a slice for later assertions.
func MockCallback(err error) {
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
	}
}

var errorMessages []string

func BenchmarkProcessDataWithConcurrency_AllSuccess(b *benchmark.B) {
	data := make([][]byte, b.N)
	callbacks := make([]func(error), b.N)

	for i := 0; i < b.N; i++ {
		data[i] = []byte{1, 2, 3}
		callbacks[i] = func(err error) {
			if err != nil {
				panic(err)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processDataWithConcurrency(data[i:i+1], callbacks[i:i+1])
	}
}

func BenchmarkProcessDataWithConcurrency_AllErrors(b *benchmark.B) {
	data := make([][]byte, b.N)
	callbacks := make([]func(error), b.N)

	for i := 0; i < b.N; i++ {
		data[i] = []byte{}
		callbacks[i] = func(err error) {
			if err != nil {
				panic(err)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processDataWithConcurrency(data[i:i+1], callbacks[i:i+1])
	}
}

func BenchmarkProcessDataWithConcurrency_Mixed(b *benchmark.B) {
	data := make([][]byte, b.N)
	callbacks := make([]func(error), b.N)

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			data[i] = []byte{}
		} else {
			data[i] = []byte{1, 2, 3}
		}
		callbacks[i] = func(err error) {
			if err != nil {
				panic(err)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processDataWithConcurrency(data[i:i+1], callbacks[i:i+1])
	}
}
