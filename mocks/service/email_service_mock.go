package mocks

import (
	"github.com/stretchr/testify/mock"
)

type EmailServiceMock struct {
	mock.Mock
}

func NewEmailServiceMock() *EmailServiceMock {
	return &EmailServiceMock{}
}

func (e *EmailServiceMock) SendEmail(toEmail string, subject string, templateFile string, data interface{}) error {
	args := e.Called(toEmail, subject, templateFile, data)
	return args.Error(0)
}
