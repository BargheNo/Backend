package usecase

type OTPService interface {
	GenerateOTP() (string, int, error)
	VerifyOTP(redisKey, otp string) error
}
