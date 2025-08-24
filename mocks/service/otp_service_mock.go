package mocks

import (
	"github.com/stretchr/testify/mock"
)

type OTPServiceMock struct {
	mock.Mock
}

func NewOTPServiceMock() *OTPServiceMock {
	return &OTPServiceMock{}
}

func (o *OTPServiceMock) GenerateOTP() (string, int, error) {
	args := o.Called()
	return args.String(0), args.Int(1), args.Error(2)
}

func (o *OTPServiceMock) VerifyOTP(redisKey, otp string) error {
	args := o.Called(redisKey, otp)
	return args.Error(0)
}
