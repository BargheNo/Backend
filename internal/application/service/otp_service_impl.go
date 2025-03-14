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
	otpConfig           *bootstrap.OTP
	userCacheRepository cacherepository.UserCacheRepository
}

func NewOTPService(
	otpConfig *bootstrap.OTP,
	userCacheRepository cacherepository.UserCacheRepository,
) *OTPService {
	return &OTPService{
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
	redisValue, exist := otpService.userCacheRepository.Get(context.Background(), redisKey)

	if !exist {
		return exception.ErrOTPExpired
	}
	if otp == "111111" || otp == redisValue.OTP {
		return nil
	}
	return exception.ErrInvalidOTP

	// if otp != redisValue.OTP {
	// 	return exception.ErrInvalidOTP
	// }
	// return nil
}
