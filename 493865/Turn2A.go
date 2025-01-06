package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Benchmark interface {
	Run() (time.Duration, error)
}

type Test1 struct{}

func (t *Test1) Run() (time.Duration, error) {
	start := time.Now()
	// Benchmark logic here
	end := time.Now()
	return end.Sub(start), nil
}

func calculateMean(durations []time.Duration) float64 {
	sum := float64(0)
	for _, duration := range durations {
		sum += float64(duration)
	}
	return sum / float64(len(durations))
}

func calculateStandardDeviation(durations []time.Duration, mean float64) float64 {
	var variance float64
	for _, duration := range durations {
		diff := float64(duration) - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(len(durations)))
}

func calculateZScore(duration time.Duration, mean, stdDev float64) float64 {
	return (float64(duration) - mean) / stdDev
}

func detectOutliers(durations []time.Duration, threshold float64) []time.Duration {
	mean := calculateMean(durations)
	stdDev := calculateStandardDeviation(durations, mean)
	outliers := []time.Duration{}
	for _, duration := range durations {
		zScore := calculateZScore(duration, mean, stdDev)
		if math.Abs(zScore) > threshold {
			outliers = append(outliers, duration)
			fmt.Println("Detected outlier:", duration, ", Z-score:", zScore)
		}
	}
	return outliers
}

func main() {
	var wg sync.WaitGroup
	tests := []Benchmark{&Test1{}}

	for _, test := range tests {
		wg.Add(1)
		go func(test Benchmark) {
			defer wg.Done()
			durations := []time.Duration{}
			for i := 0; i < 10; i++ { // Run test 10 times
				duration, err := test.Run()
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				durations = append(durations, duration)
			}

			// Detect and log outliers
			detectOutliers(durations, 2.0)

			// Calculate statistics
			mean := calculateMean(durations)
			fmt.Println("Mean duration:", mean)
		}(test)
	}

	wg.Wait()
}
