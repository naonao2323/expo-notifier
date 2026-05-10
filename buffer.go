package exponotifier

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/naonao2323/expo-notifier/internal/ring"
	"github.com/naonao2323/expo-notifier/internal/scheduler"
)

// Default threshold values used when BufferSetting fields are zero.
const (
	DefaultDelayThreshold    = time.Second
	DefaultCountThreshold    = 10
	DefaultByteThreshold     = 1 << 20 // 1MB
	DefaultBufferedByteLimit = 1 << 30 // 1GB
	DefaultWorkerLimit       = 1
)

// ErrOversizedItem is returned by Add when the item size exceeds BufferedByteLimit.
var ErrOversizedItem = errors.New("buffer: item size exceeds byte limit")

type bundle struct {
	msgs  []PushMessage
	bytes int
}

// Buffer batches PushMessages and dispatches them to a handler when count,
// byte, or delay thresholds are exceeded.
type Buffer struct {
	BufferSetting

	handler func([]PushMessage)

	sem        *semaphore.Weighted
	mu         sync.Mutex
	sched      *scheduler.Scheduler[bundle]
	cur        []PushMessage
	curBytes   int
	flushTimer *time.Timer
}

// NewBuffer creates a Buffer that calls handler with each flushed batch.
func NewBuffer(handler func([]PushMessage), setting BufferSetting) *Buffer {
	b := &Buffer{
		BufferSetting: setting,
		handler:       handler,
	}
	if b.DelayThreshold == 0 {
		b.DelayThreshold = DefaultDelayThreshold
	}
	if b.CountThreshold == 0 {
		b.CountThreshold = DefaultCountThreshold
	}
	if b.ByteThreshold == 0 {
		b.ByteThreshold = DefaultByteThreshold
	}
	if b.BufferedByteLimit == 0 {
		b.BufferedByteLimit = DefaultBufferedByteLimit
	}
	if b.WorkerLimit == 0 {
		b.WorkerLimit = DefaultWorkerLimit
	}
	b.sem = semaphore.NewWeighted(int64(b.BufferedByteLimit))
	r := ring.New[bundle](ring.Cap(b.BufferedByteLimit, b.ByteThreshold))
	b.sched = scheduler.New(&r, b.WorkerLimit, func(bdl bundle) {
		b.sem.Release(int64(bdl.bytes))
		b.handler(bdl.msgs)
	})
	return b
}

// Add enqueues msg. Blocks until byte capacity is available or ctx is done.
// Returns ErrOversizedItem if size exceeds ByteLimit, or ctx.Err() if cancelled.
func (b *Buffer) Add(ctx context.Context, msg PushMessage, size int) error {
	if size > b.BufferedByteLimit {
		return ErrOversizedItem
	}
	if err := b.sem.Acquire(ctx, int64(size)); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cur = append(b.cur, msg)
	b.curBytes += size

	if len(b.cur) >= b.CountThreshold || b.curBytes >= b.ByteThreshold {
		b.enqueueCur()
	}
	if b.flushTimer == nil {
		b.flushTimer = time.AfterFunc(b.DelayThreshold, b.tryFlush)
	}
	return nil
}

// Flush dispatches any buffered messages and waits for all handlers to return.
func (b *Buffer) Flush() {
	b.mu.Lock()
	if b.flushTimer != nil {
		b.flushTimer.Stop()
		b.flushTimer = nil
	}
	b.enqueueCur()
	b.mu.Unlock()
	b.sched.Wait()
}

func (b *Buffer) tryFlush() {
	b.mu.Lock()
	b.flushTimer = nil
	b.enqueueCur()
	b.mu.Unlock()
}

// enqueueCur moves cur into the scheduler as a bundle. Requires b.mu locked.
func (b *Buffer) enqueueCur() {
	if len(b.cur) == 0 {
		return
	}
	bdl := bundle{msgs: b.cur, bytes: b.curBytes}
	b.cur = nil
	b.curBytes = 0
	b.sched.Enqueue(bdl)
}
