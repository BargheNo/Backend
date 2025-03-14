package service

import "github.com/golang-jwt/jwt/v5"

type JWTService interface {
	GenerateToken(userID uint) (string, string)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}
