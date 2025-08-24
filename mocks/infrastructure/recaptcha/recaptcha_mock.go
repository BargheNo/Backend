package mocks

import (
	"github.com/stretchr/testify/mock"
)

type RecaptchaMock struct {
	mock.Mock
}

func NewRecaptchaMock() *RecaptchaMock {
	return &RecaptchaMock{}
}

func (r *RecaptchaMock) Verify(token string) error {
	args := r.Called(token)
	return args.Error(0)
}
