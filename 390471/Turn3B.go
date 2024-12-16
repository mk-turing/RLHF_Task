package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

const (
	minPoolSize     = 1
	maxPoolSize     = 10
	initialPoolSize = 5
	retryAttempts   = 3
	backoffDuration = time.Millisecond * 100
	maxBackoff      = time.Second * 1
)

var (
	poolSize          = int32(initialPoolSize)
	currentWorkerLoad = int32(0)
)

// workerPool manages a pool of worker goroutines
type workerPool struct {
	workers      chan int
	errors       chan error
	tasks        chan func()
	shutdown     chan struct{}
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	timeout      time.Duration
	retryBackoff func(attempt int) time.Duration
}

func newWorkerPool(timeout time.Duration) *workerPool {
	ctx, cancel := context.WithCancel(context.Background())
	retryBackoff := func(attempt int) time.Duration {
		duration := backoffDuration * time.Duration(1<<uint(attempt))
		if duration > maxBackoff {
			duration = maxBackoff
		}
		return duration
	}
	pool := &workerPool{
		workers:      make(chan int, int(poolSize)),
		errors:       make(chan error),
		tasks:        make(chan func()),
		shutdown:     make(chan struct{}),
		wg:           sync.WaitGroup{},
		ctx:          ctx,
		cancel:       cancel,
		timeout:      timeout,
		retryBackoff: retryBackoff,
	}
	go pool.start()
	return pool
}

func (p *workerPool) start() {
	for i := 0; i < int(poolSize); i++ {
		p.workers <- 1
	}
	defer close(p.errors)

	for {
		select {
		case <-p.shutdown:
			return
		case task := <-p.tasks:
			if err := p.runTaskWithRetry(task); err != nil {
				p.errors <- err
			}
		case <-time.After(1 * time.Second):
			p.adjustPoolSize()
		}
	}
}

func (p *workerPool) adjustPoolSize() {
	newSize := atomic.LoadInt32(&currentWorkerLoad)
	if newSize > int32(maxPoolSize) {
		return // Cap pool size
	}
	if newSize < int32(minPoolSize) {
		newSize = int32(minPoolSize)
	}

	currentSize := int(atomic.LoadInt32(&poolSize))
	diff := int(newSize) - currentSize

	if diff > 0 {
		for i := 0; i < diff; i++ {
			p.workers <- 1
		}
		fmt.Println("Increasing worker pool to:", newSize)
	} else if diff < 0 {
		for i := 0; i < -diff; i++ {
			<-p.workers
		}
		fmt.Println("Decreasing worker pool to:", newSize)
	}

	atomic.StoreInt32(&poolSize, newSize)
}

func (p *workerPool) runTaskWithRetry(task func()) error {
	for attempt := 0; attempt < retryAttempts; attempt++ {
		err := func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			atomic.AddInt32(&currentWorkerLoad, 1)
			defer atomic.AddInt32(&currentWorkerLoad, -1)
			defer p.wg.Done()
			return nil
		}(p.ctx)

		if err == nil {
			return nil
		}

		time.Sleep(p.retryBackoff(attempt))
		fmt.Println("Attempt", attempt+1, "failed:", err)
	}
	return fmt.Errorf("max retry attempts exceeded")
}

func (p *workerPool) Submit(task func()) {
	p.wg.Add(1)
	select {
	case p.workers <- 1:
		go func() {
			p.tasks <- task
		}()
	default:
		fmt.Println("Worker pool full, dropping task")
	}
}

func (p *workerPool) Shutdown() {
	p.cancel()
	p.shutdown <- struct{}{}
	p.wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	counter := &Counter{}
	numGoroutines := 100
	numIncrements := 1000

	// Increment counter concurrently
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				counter.Inc()
			}
		}()
	}
	wg.Wait()
	fmt.Println("Final counter value:", counter.Get())

	// Create worker pool
	pool := newWorkerPool(2 * time.Second)
	defer pool.Shutdown()

	for i := 0; i < 5; i++ {
		pool.Submit(func() {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			// Simulate a failure randomly
			if rand.Intn(2) == 0 {
			}
		})
	}

	// Handle logged errors
	go func() {
		for err := range pool.errors {
			if err != nil {
				fmt.Println("Error occurred:", err)
			}
		}
	}()

	// Wait for all workers
	pool.wg.Wait()
	fmt.Println("All workers completed successfully.")
}
