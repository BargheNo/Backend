package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JWTServiceMock struct {
	mock.Mock
}

func NewJWTServiceMock() *JWTServiceMock {
	return &JWTServiceMock{}
}

func (j *JWTServiceMock) GenerateToken(userID uint) (string, string, error) {
	args := j.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

func (j *JWTServiceMock) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	args := j.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
