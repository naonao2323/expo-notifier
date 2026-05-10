package exponotifier

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrMalformedToken is returned when a push token does not start with "ExponentPushToken".
var ErrMalformedToken = errors.New("token should start with ExponentPushToken")

const (
	// SuccessStatus is the status value returned by Expo on a successful notification.
	SuccessStatus = "ok"
	// ErrorDeviceNotRegistered indicates the push token is no longer valid.
	ErrorDeviceNotRegistered = "DeviceNotRegistered"
	// ErrorMessageTooBig indicates the notification payload exceeded 4096 bytes.
	ErrorMessageTooBig = "MessageTooBig"
	// ErrorMessageRateExceeded indicates messages are being sent too frequently.
	ErrorMessageRateExceeded = "MessageRateExceeded"
)

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

func (r *PushResponse) isSuccess() bool {
	return r.Status == SuccessStatus
}

// ValidateResponse returns a typed error if the push response indicates a failure.
func (r *PushResponse) ValidateResponse() error {
	if r.isSuccess() {
		return nil
	}
	err := &PushResponseError{Response: r}
	if r.Details != nil {
		switch r.Details["error"] {
		case ErrorDeviceNotRegistered:
			return &DeviceNotRegisteredError{PushResponseError: *err}
		case ErrorMessageTooBig:
			return &MessageTooBigError{PushResponseError: *err}
		case ErrorMessageRateExceeded:
			return &MessageRateExceededError{PushResponseError: *err}
		}
	}
	return err
}

// PushResponseError is the base error type for push notification response failures.
type PushResponseError struct {
	Response *PushResponse
}

func (e *PushResponseError) Error() string {
	if e.Response != nil {
		return fmt.Sprintf("push response error: %s", e.Response.Message)
	}
	return "unknown push response error"
}

// DeviceNotRegisteredError is returned when the push token is invalid or unregistered.
type DeviceNotRegisteredError struct{ PushResponseError }

// MessageTooBigError is returned when the notification payload exceeds the 4096-byte limit.
type MessageTooBigError struct{ PushResponseError }

// MessageRateExceededError is returned when messages are sent too frequently to a device.
type MessageRateExceededError struct{ PushResponseError }

// PushServerError is returned when the Expo push server returns an unexpected error response.
type PushServerError struct {
	Message      string
	Response     *http.Response
	ResponseData *response
	Err          error
}

func (e *PushServerError) Error() string { return e.Message }

func (e *PushServerError) Unwrap() error { return e.Err }

// Is reports whether the target is a PushServerError.
func (e *PushServerError) Is(target error) bool {
	_, ok := target.(*PushServerError)
	return ok
}

// As sets the target to this PushServerError if the target type matches.
func (e *PushServerError) As(target any) bool {
	t, ok := target.(**PushServerError)
	if !ok {
		return false
	}
	*t = e
	return true
}
