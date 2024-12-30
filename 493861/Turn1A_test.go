package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"syscall"
	"time"
)

// Function to be timed
func functionToTime() {
	for i := 0; i < 1000000; i++ {
		rand.Intn(100)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(1) // Ensure single-threaded execution for consistent timing

	// Time using time.Now()
	fmt.Println("Timing using time.Now():")
	benchmark("time.Now()", 100, func() {
		start := time.Now()
		functionToTime()
		end := time.Now()
		fmt.Printf("Duration: %v\n", end.Sub(start))
	})

	// Time using time.Monotonic()
	fmt.Println("Timing using time.Monotonic():")
	benchmark("time.Monotonic()", 100, func() {
		start := time.Monotonic()
		functionToTime()
		end := time.Monotonic()
		fmt.Printf("Duration: %v\n", end.Sub(start))
	})

	// Time using time.UnixNano()
	fmt.Println("Timing using time.UnixNano():")
	benchmark("time.UnixNano()", 100, func() {
		start := time.UnixNano()
		functionToTime()
		end := time.UnixNano()
		fmt.Printf("Duration: %v\n", end.Sub(start))
	})

	// Time using syscall.ClockGettime(syscall.CLOCK_REALTIME, &ts)
	fmt.Println("Timing using syscall.ClockGettime(syscall.CLOCK_REALTIME):")
	benchmark("syscall.ClockGettime(syscall.CLOCK_REALTIME)", 100, func() {
		start := syscall.ClockGettime(syscall.CLOCK_REALTIME, &syscall.Timespec{})
		functionToTime()
		end := syscall.ClockGettime(syscall.CLOCK_REALTIME, &syscall.Timespec{})
		fmt.Printf("Duration: %v\n", end.Sub(start))
	})

	// Time using syscall.ClockGettime(syscall.CLOCK_MONOTONIC, &ts)
	fmt.Println("Timing using syscall.ClockGettime(syscall.CLOCK_MONOTONIC):")
	benchmark("syscall.ClockGettime(syscall.CLOCK_MONOTONIC)", 100, func() {
		start := syscall.ClockGettime(syscall.CLOCK_MONOTONIC, &syscall.Timespec{})
		functionToTime()
		end := syscall.ClockGettime(syscall.CLOCK_MONOTONIC, &syscall.Timespec{})
		fmt.Printf("Duration: %v\n", end.Sub(start))
	})
}

func benchmark(name string, iterations int, block func()) {
	for i := 0; i < iterations; i++ {
		block()
	}
}
