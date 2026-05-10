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

// sendRequest sends messages to the Expo push API and returns responses or an error.
func (c *Client) sendRequest(messages []PushMessage) ([]PushResponse, error) {
	if len(messages) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(messages)
	if err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to marshal messages: %w", err)}
	}
	req, err := c.newPushRequest(body)
	if err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to build request: %w", err)}
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{Err: fmt.Errorf("request failed: %w", err)}
	}
	defer func() { _ = resp.Body.Close() }()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if len(r.Errors) > 0 {
		return nil, &PushServerError{Message: r.Errors.Error(), Err: fmt.Errorf("push server error: %w", r.Errors)}
	}
	for i := range r.Data {
		if i < len(messages) {
			r.Data[i].PushMessage = messages[i]
		}
	}
	return r.Data, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	return &PushServerError{
		Message:  fmt.Sprintf("invalid response (%d %s)", resp.StatusCode, resp.Status),
		Response: resp,
		Err:      fmt.Errorf("%d %s", resp.StatusCode, resp.Status),
	}
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
