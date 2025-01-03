package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"
)

var (
	filePath   = flag.String("file", "/etc/passwd", "Path to the file to stat")
	iterations = flag.Int("iter", 10000, "Number of times to stat the file")
	latency    = flag.Duration("latency", 0, "Simulated system call latency in nanoseconds")
)

func mockStat(path string) (syscall.Stat_t, error) {
	time.Sleep(*latency)
	return os.Stat(path)
}

func main() {
	flag.Parse()

	start := time.Now()

	for i := 0; i < *iterations; i++ {
		_, err := mockStat(*filePath)
		if err != nil {
			fmt.Println("Error statting file:", err)
			os.Exit(1)
		}
	}

	duration := time.Since(start)
	fmt.Printf("Time taken: %s\n", duration)
}
