package mocks

import "github.com/stretchr/testify/mock"

type OtpServiceMock struct {
	mock.Mock
}

func NewOtpServiceMock() *OtpServiceMock {
	return &OtpServiceMock{}
}

func (s *OtpServiceMock) GenerateOTP() (string, int) {
	args := s.Called()
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
	return args.String(0), args.Int(1)
}

func (s *OtpServiceMock) VerifyOTP(redisKey, otp string) error {
	args := s.Called(redisKey, otp)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
	return args.Error(0)
}
