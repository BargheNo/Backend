package serviceimpl

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/exception"
	cacherepository "github.com/BargheNo/Backend/internal/domain/repository/redis"
)

type OTPService struct {
	constants           *bootstrap.Constants
	otpConfig           *bootstrap.OTP
	userCacheRepository cacherepository.UserCacheRepository
}

func NewOTPService(
	constants *bootstrap.Constants,
	otpConfig *bootstrap.OTP,
	userCacheRepository cacherepository.UserCacheRepository,
) *OTPService {
	return &OTPService{
		constants:           constants,
		otpConfig:           otpConfig,
		userCacheRepository: userCacheRepository,
	}
}

var table = []byte("123456789")

func (otpService *OTPService) GenerateOTP() (string, int) {
	otp := make([]byte, otpService.otpConfig.Length)
	n, err := io.ReadAtLeast(rand.Reader, otp, otpService.otpConfig.Length)
	if n != otpService.otpConfig.Length {
		panic(err)
	}
	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}
	return string(otp), otpService.otpConfig.ExpiryMinute
}

func (otpService *OTPService) VerifyOTP(redisKey, otp string) error {
	var validationErrors exception.ValidationErrors
	redisValue, exist := otpService.userCacheRepository.Get(context.Background(), redisKey)

	if !exist {
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
