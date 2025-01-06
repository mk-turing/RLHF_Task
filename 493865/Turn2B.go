package main
import (
	"fmt"
	"sort"
	"time"
	"math"
	"sync"
)
type Benchmark interface {
	Run() (time.Duration, error)
}
type Test1 struct {}
func (t *Test1) Run() (time.Duration, error) {
	start := time.Now()
	// Benchmark logic here
	end := time.Now()
	return end.Sub(start), nil
}
// Function to calculate z-score of a given slice of float64
func calculateZScore(data []float64) []float64 {
	zScores := make([]float64, len(data))
	mean := calculateMean(data)
	stdDev := calculateStdDev(data, mean)
	for i, value := range data {
		zScores[i] = (value - mean) / stdDev
	}
	return zScores
}

// Function to calculate mean of a given slice of float64
func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Function to calculate standard deviation of a given slice of float64
func calculateStdDev(data []float64, mean float64) float64 {
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-mean, 2)
	}
	variance /= float64(len(data) - 1)
	return math.Sqrt(variance)
}

func main() {
	var wg sync.WaitGroup
	numIterations := 100
	tests := []Benchmark{&Test1{}}
	for _, test := range tests {
		wg.Add(1)
		go func(test Benchmark) {
			defer wg.Done()
			var durations []time.Duration
			for i := 0; i < numIterations; i++ {
				duration, err := test.Run()
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				durations = append(durations, duration)
			}
			// Convert durations to milliseconds for analysis
			var floatDurations []float64
			for _, duration := range durations {
				floatDurations = append(floatDurations, float64(duration/time.Millisecond))
			}

			// Calculate z-scores
			zScores := calculateZScore(floatDurations)

			// Identify outliers based on z-scores (3-sigma rule)
			threshold := 3.0
			outliers := make([]float64, 0)
			for _, zScore := range zScores {
				if math.Abs(zScore) > threshold {
					outliers = append(outliers, zScore)
				}
			}
			sort.Float64s(outliers)
			if len(outliers) > 0 {
				fmt.Println("Detected outliers:")
				for _, outlier := range outliers {
					fmt.Printf("Z-Score: %0.2f\n", outlier)
				}
			}

			// Calculate median and mean ignoring outliers
			var filteredDurations []float64
			for i, value := range floatDurations {
				if math.Abs(zScores[i]) <= threshold {
					filteredDurations = append(filteredDurations, value)
				}
			}
			mean := calculateMean(filteredDurations)
			median := median(filteredDurations)
			fmt.Println("Performance Metrics (ignoring outliers):")
			fmt.Printf("Mean: %.2f ms\n", mean)
			fmt.Printf("Median: %.2f ms\n", median)
		}(test)
	}
	wg.Wait()
}
//Function to calculate median
func median(numbers []float64) float64 {
	n := len(numbers)
	if n == 0 {<