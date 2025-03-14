package service

type OTPService interface {
	GenerateOTP() (string, int)
	VerifyOTP(redisKey, otp string) error
}
