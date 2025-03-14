package service

type SMSService interface {
	SendOTP(receptor, token string)
}
