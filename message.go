package exponotifier

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// DefaultPriority is the standard delivery priority for push messages.
	DefaultPriority = "default"
	// NormalPriority is the normal delivery priority for push messages.
	NormalPriority = "normal"
	// HighPriority is the high delivery priority for push messages.
	HighPriority = "high"
)

// ExponentPushToken is a validated Expo push token.
type ExponentPushToken string

// ErrMalformedToken is returned when a push token does not start with "ExponentPushToken".
var ErrMalformedToken = errors.New("token should start with ExponentPushToken")

// NewExponentPushToken returns a validated ExponentPushToken or an error if the token is malformed.
func NewExponentPushToken(token string) (ExponentPushToken, error) {
	if !strings.HasPrefix(token, "ExponentPushToken") {
		return "", ErrMalformedToken
	}
	return ExponentPushToken(token), nil
}

// PushMessage describes a push notification request.
type PushMessage struct {
	To         ExponentPushToken `json:"to"`
	Title      string            `json:"title,omitempty"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
	Sound      string            `json:"sound,omitempty"`
	TTLSeconds int               `json:"ttl,omitempty"`
	Expiration int64             `json:"expiration,omitempty"`
	Priority   string            `json:"priority,omitempty"`
	Badge      int               `json:"badge,omitempty"`
	ChannelID  string            `json:"channelId,omitempty"`
}

type response struct {
	Data   []PushResponse `json:"data"`
	Errors apiErrors      `json:"errors"`
}

// PushResponse holds the result of a single push notification request.
type PushResponse struct {
	PushMessage PushMessage
	ID          string            `json:"id"`
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details"`
}

type apiErrors []map[string]string

func (e apiErrors) Error() string {
	msgs := make([]string, len(e))
	for i, m := range e {
		msgs[i] = fmt.Sprintf("%v", m)
	}
	return strings.Join(msgs, "\n")
}

func (e apiErrors) Unwrap() []error {
	errs := make([]error, len(e))
	for i, m := range e {
		errs[i] = fmt.Errorf("%v", m)
	}
	return errs
}
