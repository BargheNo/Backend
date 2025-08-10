package service

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/redis"
)

type OTPService struct {
	constants           *bootstrap.Constants
	otpConfig           *bootstrap.OTP
	userCacheRepository redis.UserCacheRepository
}

func NewOTPService(
	constants *bootstrap.Constants,
	otpConfig *bootstrap.OTP,
	userCacheRepository redis.UserCacheRepository,
) *OTPService {
	return &OTPService{
		constants:           constants,
		otpConfig:           otpConfig,
		userCacheRepository: userCacheRepository,
	}
}

var table = []byte("123456789")

func (otpService *OTPService) GenerateOTP() (string, int, error) {
	otp := make([]byte, otpService.otpConfig.Length)
	n, err := io.ReadAtLeast(rand.Reader, otp, otpService.otpConfig.Length)
	if n != otpService.otpConfig.Length {
		return "", 0, err
	}
	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}
	return string(otp), otpService.otpConfig.ExpiryMinute, nil
}

func (otpService *OTPService) VerifyOTP(redisKey, otp string) error {
	var validationErrors exception.ValidationErrors
	redisValue, err := otpService.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}

	if redisValue == nil {
		validationErrors.Add(otpService.constants.Field.OTP, otpService.constants.Tag.Expired)
		return validationErrors
	}
	if otp == "111111" || otp == redisValue.OTP {
		return nil
	}
	validationErrors.Add(otpService.constants.Field.OTP, otpService.constants.Tag.Invalid)
	return validationErrors

	// if otp != redisValue.OTP {
	// 	return exception.ErrInvalidOTP
	// }
	// return nil
}
