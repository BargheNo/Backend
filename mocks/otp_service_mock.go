package mocks

import "github.com/stretchr/testify/mock"

type OtpServiceMock struct {
	mock.Mock
}

func NewOtpServiceMock() *OtpServiceMock {
	return &OtpServiceMock{}
}

func (s *OtpServiceMock) GenerateOTP() (string, int, error) {
	args := s.Called()
	return args.String(0), args.Int(1), args.Error(2)
}

func (s *OtpServiceMock) VerifyOTP(redisKey, otp string) error {
	args := s.Called(redisKey, otp)
	return args.Error(0)
}
