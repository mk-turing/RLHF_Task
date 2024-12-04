package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// PanicsAnnotatedError implement PolicyError for errors that should panic
type PanicsAnnotatedError struct {
	err error
}

func (e *PanicsAnnotatedError) Error() string {
	return e.err.Error()
}

func (e *PanicsAnnotatedError) GetPolicy() string {
	return "panic"
}

// LogsAnnotatedError implement PolicyError for errors that should be logged
type LogsAnnotatedError struct {
	err error
}

func (e *LogsAnnotatedError) Error() string {
	return e.err.Error()
}

func (e *LogsAnnotatedError) GetPolicy() string {
	return "log"
}

// ReturnsAnnotatedError implement PolicyError for errors that should be returned
type ReturnsAnnotatedError struct {
	err error
}

func (e *ReturnsAnnotatedError) Error() string {
	return e.err.Error()
}

func (e *ReturnsAnnotatedError) GetPolicy() string {
	return "return"
}

// RetryAnnotatedError extends PolicyError for errors that should be retried
type RetryAnnotatedError struct {
	err      error
	maxRetry int
	retryCnt int
}

func (e *RetryAnnotatedError) Error() string {
	return fmt.Sprintf("retry %d/%d: %v", e.retryCnt, e.maxRetry, e.err)
}

func (e *RetryAnnotatedError) GetPolicy() string {
	return "retry"
}

func (e *RetryAnnotatedError) CanRetry() bool {
	return e.retryCnt < e.maxRetry
}

func retryHandler(f func() error, maxRetry int, retryInterval time.Duration) {
	err := f()
	for retry := range retryLoop(err, maxRetry, retryInterval) {
		if retry {
			time.Sleep(retryInterval)
		}
		if err == nil {
			return
		}
		handleError(err)
	}
	handleError(err)
}

func retryLoop(err error, maxRetry int, retryInterval time.Duration) <-chan bool {
	ch := make(chan bool, maxRetry+1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for attempt := 0; attempt <= maxRetry; attempt++ {
			select {
			case ch <- attempt < maxRetry && (err != nil):
				if retryErr, ok := err.(*RetryAnnotatedError); ok {
					retryErr.retryCnt = attempt + 1
				}
			case <-time.After(retryInterval):
			}
			if err == nil {
				ch <- false
				return
			}
			if attempt < maxRetry {
				err = retry(err)
			}
		}
		ch <- false
	}()
	return ch
}

func retry(err error) error {
	return err
}

func externalLogger(msg string) {
	log.Println("External Log:", msg)
}

func handleError(err error) {
	switch err := err.(type) {
	case *PanicsAnnotatedError:
		panic(err.err)
	case *LogsAnnotatedError, *RetryAnnotatedError:
		externalLogger(err.Error())
	case *ReturnsAnnotatedError:
		fmt.Println("Function returned error:", err.err)
	default:
		externalLogger("Unknown error type: " + err.Error())
	}
}

func readFromFile(file string, retries int, interval time.Duration) ([]byte, error) {
	for retries >= 0 {
		data, err := os.ReadFile(file)
		if err == nil {
			return data, nil
		}
		if retries == 0 {
			return nil, &ReturnsAnnotatedError{err: err}
		}
		retries--
		time.Sleep(interval)
	}
	return nil, nil
}

// Function that "panics" on error
func divide(a, b int) (err error) {
	if b == 0 {
		return errors.New("divide by zero")
	}
	result := a / b
	fmt.Println(result)
	return nil
}

// Function that "logs" on error
func writeToFile(file string) error {
	if err := os.WriteFile(file, nil, 0644); err != nil {
		return &LogsAnnotatedError{err: err}
	}
	return nil
}

func divideByZero(retries int, interval time.Duration) error {
	if err := divide(1, 0); err != nil {
		return &RetryAnnotatedError{
			err:      err,
			maxRetry: retries,
			retryCnt: 0,
		}
	}
	return nil
}

func main() {
	retryHandler(func() error {
		return writeToFile("non_existent_file.txt")
	}, 2, 1*time.Second)

	_, err := readFromFile("non_existent_file.txt", 3, 1*time.Second)
	if err != nil {
		fmt.Println("Reading failed:", err)
	}

	err = divideByZero(1, 1*time.Second)
	handleError(err)

	num, _ := strconv.Atoi("xyz")
	if num == 0 {
		return
	}
}
