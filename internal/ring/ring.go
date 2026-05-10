// Package ring provides a fixed-capacity generic circular buffer.
package ring

// Ring is a fixed-capacity circular buffer. Not safe for concurrent use.
// len(buf) must be a power of 2 so that (idx+1)&(len-1) wraps correctly.
type Ring[T any] struct {
	buf  []T
	head int
	tail int
	len  int
}

// New returns a Ring backed by a buffer of the next power-of-2 size >= n.
func New[T any](n int) Ring[T] {
	return Ring[T]{buf: make([]T, roundUpPow2(n))}
}

// Push appends v to the tail. If the buffer is full, the oldest item is overwritten.
func (r *Ring[T]) Push(v T) {
	r.buf[r.tail] = v
	r.tail = (r.tail + 1) & r.mask()
	if r.len < len(r.buf) {
		r.len++
	}
}

// Pop removes and returns the item at the head.
func (r *Ring[T]) Pop() T {
	v := r.buf[r.head]
	var zero T
	r.buf[r.head] = zero
	r.head = (r.head + 1) & r.mask()
	r.len--
	return v
}

// Len returns the number of items currently in the buffer.
func (r *Ring[T]) Len() int { return r.len }

func (r Ring[T]) mask() int {
	return (len(r.buf) - 1)
}

func roundUpPow2(n int) int {
	if n < 1 {
		n = 1
	}
	c := 1
	for c < n {
		c <<= 1
	}
	return c
}
