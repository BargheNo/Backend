package exception

type ForbiddenError struct {
	Resource string
	Message  string
}

func (e ForbiddenError) Error() string {
	if e.Message != "" {
		return "Forbidden: " + e.Message
	}
	return "Forbidden: Access to " + e.Resource + " is not allowed"
}
