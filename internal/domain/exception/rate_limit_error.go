package exception

type RateLimitError struct {
	Message string
}

func (e RateLimitError) Error() string {
	return e.Message
}
