package main

import (
	"math/rand"
	"strconv"
	"testing"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// generateRandomString generates a random string of specified length.
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func addPrefix(strings []string, prefix string) []string {
	var prefixedStrings []string
	for _, str := range strings {
		prefixedStrings = append(prefixedStrings, prefix+str)
	}
	return prefixedStrings
}

// addPrefix benchmarks the performance of the addPrefix function for different prefix sizes and dataset sizes.
func BenchmarkAddPrefix(b *testing.B) {
	testCases := []struct {
		prefixSize  int
		datasetSize int
	}{
		{1, 10000},
		{5, 10000},
		{10, 10000},
		{1, 100000},
		{5, 100000},
		{10, 100000},
	}

	for _, testCase := range testCases {
		prefix := generateRandomString(testCase.prefixSize)
		dataset := make([]string, testCase.datasetSize)
		for i := range dataset {
			dataset[i] = generateRandomString(10) // 10 character long strings
		}

		b.Run(strconv.Itoa(testCase.prefixSize)+"-"+strconv.Itoa(testCase.datasetSize), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = addPrefix(dataset, prefix)
			}
		})
	}
}
