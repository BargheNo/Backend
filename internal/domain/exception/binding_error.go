package exception

type BindingError struct {
	Err error
}

func (e BindingError) Error() string {
	return e.Err.Error()
}

func NewBindingError(err error) BindingError {
	return BindingError{
		Err: err,
	}
}
