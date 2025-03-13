package service

type OTPService interface {
	GenerateOTP() (string, int)
	// VerifyOTP(user *entity.User, inputOTP string, otpFieldError string, expiredTokenTagError string, invalidTokenTagError string)
}
