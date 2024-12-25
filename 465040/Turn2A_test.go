package _65040

import (
	"sync"
	"testing"
	"time"
)

func taskWithWaitGroup(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Millisecond * 50)
}

func taskWithMutex(mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	time.Sleep(time.Millisecond * 50)
}

func taskWithChannel(done chan<- struct{}) {
	time.Sleep(time.Millisecond * 50)
	done <- struct{}{}
}

func BenchmarkWaitGroup(b *testing.B) {
	wg := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go taskWithWaitGroup(wg)
	}

	for i := 0; i < b.N; i++ {
		wg.Wait()
	}
}

func BenchmarkMutex(b *testing.B) {
	mutex := &sync.Mutex{}

	for i := 0; i < b.N; i++ {
		go taskWithMutex(mutex)
	}

	for i := 0; i < b.N; i++ {
		taskWithMutex(mutex)
	}
}

func BenchmarkChannel(b *testing.B) {
	done := make(chan<- struct{}, b.N)

	for i := 0; i < b.N; i++ {
		go taskWithChannel(done)
	}
}
