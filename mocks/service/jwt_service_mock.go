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
	// Handle the case where the first argument might be nil or not jwt.MapClaims
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	// Try to convert to jwt.MapClaims, but handle the case where it might be a map[string]interface{}
	if claims, ok := args.Get(0).(jwt.MapClaims); ok {
		return claims, args.Error(1)
	}

	// If it's a map[string]interface{}, convert it
	if mapClaims, ok := args.Get(0).(map[string]interface{}); ok {
		return jwt.MapClaims(mapClaims), args.Error(1)
	}

	return nil, args.Error(1)
}
