package mocks

import "github.com/stretchr/testify/mock"

type EmailServiceMock struct {
	mock.Mock
}

func NewEmailServiceMock() *EmailServiceMock {
	return &EmailServiceMock{}
}

func (s *EmailServiceMock) SendEmail(toEmail string, subject string, templateFile string, data interface{}) {
	args := s.Called(toEmail, subject, templateFile, data)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}
