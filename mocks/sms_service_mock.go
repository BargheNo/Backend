package mocks

import "github.com/stretchr/testify/mock"

type SMSServiceMock struct {
	mock.Mock
}

func NewSMSServiceMock() *SMSServiceMock {
	return &SMSServiceMock{}
}

func (s *SMSServiceMock) SendOTP(receptor, token string) error {
	args := s.Called(receptor, token)
	return args.Error(0)
}
