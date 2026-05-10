// Package exponotifier provides a client for sending Expo push notifications.
package exponotifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// DefaultHost is the default Expo push notification host.
	DefaultHost = "https://exp.host"
	// DefaultBaseAPIURL is the default base path for Expo API requests.
	DefaultBaseAPIURL = "/--/api/v2"
)

// DefaultHTTPClient is the default HTTP client used for API requests.
var DefaultHTTPClient = http.DefaultClient

// Client handles raw HTTP communication with the Expo push API.
type Client struct {
	httpClient  *http.Client
	host        string
	url         string
	accessToken string
}

// sendRequest sends messages to the Expo push API; errors are silently dropped.
// Its signature matches the func([]PushMessage) handler expected by Buffer.
func (c *Client) sendRequest(messages []PushMessage) {
	if len(messages) == 0 {
		return
	}
	body, err := json.Marshal(messages)
	if err != nil {
		return
	}
	req, err := c.newPushRequest(body)
	if err != nil {
		return
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()
	var r response
	_ = json.NewDecoder(resp.Body).Decode(&r)
}

func (c *Client) newPushRequest(body []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s%s/push/send", c.host, c.url)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	return req, nil
}
