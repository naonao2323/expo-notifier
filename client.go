// Package exponotifier provides a client for sending Expo push notifications.
package exponotifier

import "net/http"

const (
	// DefaultHost is the default Expo push notification host.
	DefaultHost = "https://exp.host"
	// DefaultBaseAPIURL is the default base path for Expo API requests.
	DefaultBaseAPIURL = "/--/api/v2"
)

// DefaultHTTPClient is the default HTTP client used for API requests.
var DefaultHTTPClient = http.DefaultClient

// Notifier sends push notifications via the Expo push notification service.
type Notifier struct {
	client      *http.Client
	host        string
	url         string
	accessToken string
}

// NewNotifier creates a new Notifier with the given options.
func NewNotifier(opts ...Option) *Notifier {
	n := &Notifier{
		host:   DefaultHost,
		url:    DefaultBaseAPIURL,
		client: DefaultHTTPClient,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}
