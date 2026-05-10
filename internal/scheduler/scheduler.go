// Package scheduler provides a generic bounded worker pool that dispatches
// items from a Queue to a handler function.
package scheduler

import "sync"

// Queue is the storage backend used by Scheduler to hold pending items.
type Queue[T any] interface {
	Push(T)
	Pop() T
	Len() int
}

// Scheduler dispatches enqueued items to a handler via a bounded worker pool.
type Scheduler[T any] struct {
	queue       Queue[T]
	mu          sync.Mutex
	workerCount int
	workerLimit int
	wg          sync.WaitGroup
	handler     func(T)
}

// New creates a Scheduler that calls handler for each item popped from queue,
// spawning at most workerLimit concurrent goroutines.
func New[T any](queue Queue[T], workerLimit int, handler func(T)) *Scheduler[T] {
	return &Scheduler[T]{
		queue:       queue,
		workerLimit: workerLimit,
		handler:     handler,
	}
}

// Enqueue pushes v into the queue and spawns a worker goroutine if under the limit.
func (s *Scheduler[T]) Enqueue(v T) {
	s.mu.Lock()
	s.queue.Push(v)
	spawn := s.workerCount < s.workerLimit
	if spawn {
		s.workerCount++
		s.wg.Add(1)
	}
	s.mu.Unlock()
	if spawn {
		go s.work()
	}
}

// Wait blocks until all in-flight workers have returned.
func (s *Scheduler[T]) Wait() {
	s.wg.Wait()
}

func (s *Scheduler[T]) work() {
	for {
		s.mu.Lock()
		if s.queue.Len() == 0 {
			s.workerCount--
			s.wg.Done()
			s.mu.Unlock()
			return
		}
		v := s.queue.Pop()
		s.mu.Unlock()
		s.handler(v)
	}
}
