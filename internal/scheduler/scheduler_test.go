package scheduler_test

import (
	"slices"
	"sync"
	"testing"

	"github.com/naonao2323/expo-notifier/internal/scheduler"
)

type testQueue[T any] struct{ data []T }

func (q *testQueue[T]) Push(v T) { q.data = append(q.data, v) }
func (q *testQueue[T]) Pop() T   { v := q.data[0]; q.data = q.data[1:]; return v }
func (q *testQueue[T]) Len() int { return len(q.data) }

func TestSchedulerAllItemsProcessed(t *testing.T) {
	tests := []struct {
		name        string
		items       int
		workerLimit int
	}{
		{"single worker", 5, 1},
		{"multi worker", 9, 3},
		{"more workers than items", 3, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var got []int
			handler := func(v int) {
				mu.Lock()
				got = append(got, v)
				mu.Unlock()
			}
			q := &testQueue[int]{}
			s := scheduler.New(q, tc.workerLimit, handler)
			for i := range tc.items {
				s.Enqueue(i)
			}
			s.Wait()
			if len(got) != tc.items {
				t.Errorf("handler called %d times, want %d", len(got), tc.items)
			}
		})
	}
}

func TestSchedulerSingleWorkerFIFO(t *testing.T) {
	var mu sync.Mutex
	var got []int
	handler := func(v int) {
		mu.Lock()
		got = append(got, v)
		mu.Unlock()
	}
	q := &testQueue[int]{}
	s := scheduler.New(q, 1, handler)
	for i := 1; i <= 5; i++ {
		s.Enqueue(i)
	}
	s.Wait()
	if !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("got %v, want [1 2 3 4 5]", got)
	}
}
