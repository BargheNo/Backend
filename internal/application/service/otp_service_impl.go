package serviceimpl

import (
	"crypto/rand"
	"io"
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
)

type OTPService struct {
	otpConfig *bootstrap.OTP
}

func NewOTPService(otpConfig *bootstrap.OTP) *OTPService {
	return &OTPService{
		otpConfig: otpConfig,
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
