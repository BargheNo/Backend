package service

import (
	"errors"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/exception"
	domainJWT "github.com/BargheNo/Backend/internal/domain/jwt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	keysPath   *bootstrap.JWTKeysPath
	keyManager domainJWT.KeyManager
}

func NewJWTService(
	keyManager domainJWT.KeyManager,
	keysPath *bootstrap.JWTKeysPath,
) *JWTService {
	service := &JWTService{
		keyManager: keyManager,
		keysPath:   keysPath,
	}
	err := keyManager.LoadKeys(keysPath.PrivateKey, keysPath.PublicKey)
	if err != nil {
		panic(err)
	}

	return service
}

func (jwtService *JWTService) GenerateToken(userID uint) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"sub": userID,
		// "exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"exp": time.Now().Add(time.Second * 20).Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtService.keyManager.GetPrivateKey())
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"sub": userID,
		// "exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"exp": time.Now().Add(time.Minute * 1).Unix(),
		"iat": time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtService.keyManager.GetPrivateKey())
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (jwtService *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, exception.NewInvalidTokenError(nil)
		}
		return jwtService.keyManager.GetPublicKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, exception.NewExpiredTokenError(err)
		}
		return nil, exception.NewInvalidTokenError(err)
	}

	if !token.Valid {
		return nil, exception.NewInvalidTokenError(nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, exception.NewInvalidTokenError(nil)
	}

	return claims, nil
}
