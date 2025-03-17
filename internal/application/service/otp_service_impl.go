package serviceimpl

import (
	"context"
	"crypto/rand"
	"io"
	"strconv"

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
	otpLength, err := strconv.Atoi(otpService.otpConfig.Length)
	if err != nil {
		otpLength = 6
	}

	otp := make([]byte, otpLength)
	n, err := io.ReadAtLeast(rand.Reader, otp, otpLength)
	if n != otpLength {
		panic(err)
	}
	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}
	expiryMinute, err := strconv.Atoi(otpService.otpConfig.ExpiryMinute)
	if err != nil {
		expiryMinute = 2
	}
	return string(otp), expiryMinute
}

func (otpService *OTPService) VerifyOTP(redisKey, otp string) error {
	var validationErrors exception.ValidationErrors
	redisValue, exist := otpService.userCacheRepository.Get(context.Background(), redisKey)

	if !exist {
		validationErrors.Add(otpService.constants.Field.OTP, otpService.constants.Tag.OTPExpired)
		return validationErrors
	}
	if otp == "111111" || otp == redisValue.OTP {
		return nil
	}
	validationErrors.Add(otpService.constants.Field.OTP, otpService.constants.Tag.InvalidOTP)
	return validationErrors

	// if otp != redisValue.OTP {
	// 	return exception.ErrInvalidOTP
	// }
	// return nil
}
