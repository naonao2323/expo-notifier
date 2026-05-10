package exponotifier

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring a Notifier.
type Option func(*Notifier)

// WithHost sets the Expo push notification host.
func WithHost(host string) Option {
	return func(n *Notifier) { n.client.host = host }
}

// WithAPIURL sets the base API URL path.
func WithAPIURL(apiURL string) Option {
	return func(n *Notifier) { n.client.url = apiURL }
}

// WithAccessToken sets the Expo access token for authenticated requests.
func WithAccessToken(token string) Option {
	return func(n *Notifier) { n.client.accessToken = token }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(n *Notifier) { n.client.httpClient = client }
}

// WithBuffer enables batched sending with the given buffer options.
// Messages added via Add are accumulated and flushed as batch HTTP requests.
func WithBuffer(opts ...BufferOption) Option {
	return func(n *Notifier) {
		for _, o := range opts {
			o(&n.bufferSetting)
		}
	}
}

// BufferSetting holds throttle configuration for the Buffer.
type BufferSetting struct {
	DelayThreshold    time.Duration
	CountThreshold    int
	ByteThreshold     int // per-bundle soft trigger
	BufferedByteLimit int // total bytes allowed across ring + in-flight (semaphore capacity)
	WorkerLimit       int
}

// RetryDecision is the return value of an ErrorHandler.
type RetryDecision bool

const (
	// Retry re-enqueues the failed messages into the buffer.
	Retry RetryDecision = true
	// NoRetry discards the failed messages.
	NoRetry RetryDecision = false
)

// ErrorHandler is called when a batch send fails (network error, non-2xx status, API error)
// or when individual push responses indicate delivery failures.
// Return Retry to re-enqueue msgs into the buffer, or NoRetry to discard them.
type ErrorHandler func(msgs []PushMessage, err error) RetryDecision

// WithErrorHandler sets a callback invoked on send failures.
// Without this option, errors are silently dropped.
func WithErrorHandler(h ErrorHandler) Option {
	return func(n *Notifier) { n.errorHandler = h }
}

// BufferOption is a functional option for configuring a BufferSetting.
type BufferOption func(*BufferSetting)

// WithDelayThreshold sets the maximum time a message waits before being flushed.
func WithDelayThreshold(d time.Duration) BufferOption {
	return func(s *BufferSetting) { s.DelayThreshold = d }
}

// WithCountThreshold sets the number of messages that triggers an immediate flush.
func WithCountThreshold(n int) BufferOption {
	return func(s *BufferSetting) { s.CountThreshold = n }
}

// WithByteThreshold sets the per-bundle byte size that triggers an immediate flush.
func WithByteThreshold(n int) BufferOption {
	return func(s *BufferSetting) { s.ByteThreshold = n }
}

// WithBufferedByteLimit sets the total byte capacity across buffered and in-flight messages.
func WithBufferedByteLimit(n int) BufferOption {
	return func(s *BufferSetting) { s.BufferedByteLimit = n }
}

// WithWorkerLimit sets the maximum number of concurrent batch handler goroutines.
func WithWorkerLimit(n int) BufferOption { return func(s *BufferSetting) { s.WorkerLimit = n } }
