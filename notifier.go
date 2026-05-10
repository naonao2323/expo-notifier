package exponotifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Notify sends a single push notification and returns the push response.
func (n *Notifier) Notify(message PushMessage) (PushResponse, error) {
	if message.To == "" {
		return PushResponse{}, errors.New("no recipient")
	}
	resp, err := n.notify(message)
	if err != nil {
		return PushResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return PushResponse{}, fmt.Errorf("invalid response (%d %s)", resp.StatusCode, resp.Status)
	}
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return PushResponse{}, err
	}
	if r.Errors != nil {
		return PushResponse{}, &PushServerError{Message: "invalid server response", Response: resp, ResponseData: &r, Err: r.Errors}
	}
	if len(r.Data) == 0 {
		return PushResponse{}, &PushServerError{Message: "invalid server response", Response: resp, ResponseData: &r}
	}
	result := r.Data[0]
	result.PushMessage = message
	return result, nil
}

func (n *Notifier) notify(message PushMessage) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/push/send", n.host, n.url)
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if n.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.accessToken)
	}
	return n.client.Do(req)
}
