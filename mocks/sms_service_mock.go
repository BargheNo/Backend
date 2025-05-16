package mocks

import "github.com/stretchr/testify/mock"

type SMSServiceMock struct {
	mock.Mock
}

func NewSMSServiceMock() *SMSServiceMock {
	return &SMSServiceMock{}
}

func (s *SMSServiceMock) SendOTP(receptor, token string) {
	args := s.Called(receptor, token)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}
