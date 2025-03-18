package exception

import "fmt"

const (
	ErrorTypeRequestRateLimit       = "REQUEST_RATE_LIMIT"
	ErrorTypeConcurrentInstallLimit = "CONCURRENT_INSTALL_LIMIT"
)

type RateLimitError struct {
	Type        string
	Message     string
	Limit       int
	OriginalErr error
}

func (e RateLimitError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s: %s (limit: %d, %s)", e.Type, e.Message, e.Limit, e.OriginalErr.Error())
	}
	return fmt.Sprintf("%s: %s (limit: %d)", e.Type, e.Message, e.Limit)
}

func NewRequestRateLimitError(message string, limit int, originalErr error) *RateLimitError {
	if message == "" {
		message = "request rate limit exceeded"
	}
	return &RateLimitError{
		Type:        ErrorTypeRequestRateLimit,
		Message:     message,
		Limit:       limit,
		OriginalErr: originalErr,
	}
}

func NewConcurrentInstallLimitError(message string, limit int, originalErr error) *RateLimitError {
	if message == "" {
		message = "maximum concurrent installations limit reached"
	}
	return &RateLimitError{
		Type:        ErrorTypeConcurrentInstallLimit,
		Message:     message,
		Limit:       limit,
		OriginalErr: originalErr,
	}
}
