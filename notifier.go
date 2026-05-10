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
					data, e := json.Marshal(msg)
					if e != nil {
						continue
					}
					_ = n.buffer.Add(context.Background(), msg, len(data))
				}
			}
			return
		}
		for _, r := range responses {
			if perMsgErr := r.ValidateResponse(); perMsgErr != nil {
				if n.errorHandler([]PushMessage{r.PushMessage}, perMsgErr) == Retry {
					data, e := json.Marshal(r.PushMessage)
					if e == nil {
						_ = n.buffer.Add(context.Background(), r.PushMessage, len(data))
					}
				}
			}
		}
	}
}
