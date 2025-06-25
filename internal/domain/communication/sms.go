package communication

type SMSService interface {
	SendOTP(receptor string, token string) error
}
