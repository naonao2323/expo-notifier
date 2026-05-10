package exponotifier

import (
	"context"
)

// Notifier sends Expo push notifications via a buffer.
type Notifier struct {
	client        *Client
	buffer        *Buffer
	bufferSetting BufferSetting
	errorHandler  ErrorHandler
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
	n.buffer = NewBuffer(n.makeHandler(), n.bufferSetting)
	return n
}

// Add enqueues msg into the buffer. Blocks until byte capacity is available or ctx is done.
func (n *Notifier) Add(ctx context.Context, msg PushMessage) error {
	size := len(msg.Body) + len(msg.Title) + len(msg.Data)*64 + 128
	return n.buffer.Add(ctx, msg, size)
}

// Flush drains all buffered messages and waits for all batch sends to complete.
func (n *Notifier) Flush() {
	n.buffer.Flush()
}

func (n *Notifier) makeHandler() func([]PushMessage) {
	return func(msgs []PushMessage) {
		responses, err := n.client.sendRequest(msgs)
		if err == nil && responses == nil {
			return
		}
		if n.errorHandler == nil {
			return
		}
		if err != nil {
			if n.errorHandler(msgs, err) == Retry {
				for _, msg := range msgs {
					_ = n.Add(context.Background(), msg)
				}
			}
			return
		}
		for _, r := range responses {
			if perMsgErr := r.ValidateResponse(); perMsgErr != nil {
				if n.errorHandler([]PushMessage{r.PushMessage}, perMsgErr) == Retry {
					_ = n.Add(context.Background(), r.PushMessage)
				}
			}
		}
	}
}
