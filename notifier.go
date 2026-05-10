package exponotifier

import (
	"context"
	"encoding/json"
)

// Notifier sends Expo push notifications via a buffer.
type Notifier struct {
	client        *Client
	buffer        *Buffer
	bufferSetting BufferSetting
}

// NewNotifier creates a Notifier with the given options. A buffer is always initialized.
func NewNotifier(opts ...Option) *Notifier {
	n := &Notifier{
		client: &Client{
			httpClient: DefaultHTTPClient,
			host:       DefaultHost,
			url:        DefaultBaseAPIURL,
		},
		bufferSetting: BufferSetting{
			DelayThreshold:    DefaultDelayThreshold,
			CountThreshold:    DefaultCountThreshold,
			ByteThreshold:     DefaultByteThreshold,
			BufferedByteLimit: DefaultBufferedByteLimit,
			WorkerLimit:       DefaultWorkerLimit,
		},
	}
	for _, opt := range opts {
		opt(n)
	}
	n.buffer = NewBuffer(n.client.sendRequest, n.bufferSetting)
	return n
}

// Add enqueues msg into the buffer. Blocks until byte capacity is available or ctx is done.
func (n *Notifier) Add(ctx context.Context, msg PushMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return n.buffer.Add(ctx, msg, len(data))
}

// Flush drains all buffered messages and waits for all batch sends to complete.
func (n *Notifier) Flush() {
	n.buffer.Flush()
}
