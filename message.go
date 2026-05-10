package exponotifier

import "strings"

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
