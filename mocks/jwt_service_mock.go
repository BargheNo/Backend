package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func NewJwtServiceMock() *JwtServiceMock {
	return &JwtServiceMock{}
}

func (s *JwtServiceMock) GenerateToken(userID uint) (string, string, error) {
	args := s.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

func (s *JwtServiceMock) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	args := s.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
