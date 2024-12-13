package main

import (
	"fmt"
	"log"
	"runtime/debug"
)

type DivideError struct {
	Dividend float64
	Divisor  float64
}

func (e *DivideError) Error() string {
	return fmt.Sprintf("cannot divide %.2f by %.2f", e.Dividend, e.Divisor)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, &DivideError{Dividend: a, Divisor: b}
	}
	return a / b, nil
}

func safeExecute(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v\nstack trace:\n%s", r, debug.Stack())
		}
	}()
	fn()
	return nil
}

func main() {
	err := safeExecute(func() {
		_, err := divide(10, 0)
		if err != nil {
			log.Printf("error occurred: %v", err)
			return
		}
	})
	if err != nil {
		log.Fatalf("an error occurred: %v", err)
	}
}
