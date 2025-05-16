package mocks

import "github.com/stretchr/testify/mock"

type EmailServiceMock struct {
	mock.Mock
}

func NewEmailServiceMock() *EmailServiceMock {
	return &EmailServiceMock{}
}

func (s *EmailServiceMock) SendEmail(toEmail string, subject string, templateFile string, data interface{}) error {
	args := s.Called(toEmail, subject, templateFile, data)
	return args.Error(0)
}
