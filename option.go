package exponotifier

import "net/http"

// Option is a functional option for configuring a Notifier.
type Option func(*Notifier)

// WithHost sets the Expo push notification host.
func WithHost(host string) Option {
	return func(n *Notifier) {
		n.host = host
	}
}

// WithAPIURL sets the base API URL path.
func WithAPIURL(apiURL string) Option {
	return func(n *Notifier) {
		n.url = apiURL
	}
}

// WithAccessToken sets the Expo access token for authenticated requests.
func WithAccessToken(token string) Option {
	return func(n *Notifier) {
		n.accessToken = token
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(n *Notifier) {
		n.client = client
	}
}
