package recaptcha

type Recaptcha interface {
	Verify(token string) error
}
