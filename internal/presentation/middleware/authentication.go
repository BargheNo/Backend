package middleware

import (
	"strings"

	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	jwtService     service.JWTService
	userRepository repository.UserRepository
	db             database.Database
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	userRepository repository.UserRepository,
	db database.Database,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:      constants,
		jwtService:     jwtService,
		userRepository: userRepository,
		db:             db,
	}
}

func (am *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		unauthorizedError := exception.NewUnauthorizedError("empty auth header", nil)
		panic(unauthorizedError)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError := exception.NewUnauthorizedError("invalid token format", nil)
		panic(unauthorizedError)
	}

	tokenString := parts[1]
	if tokenString == "" {
		unauthorizedError := exception.NewUnauthorizedError("empty token", nil)
		panic(unauthorizedError)
	}

	claims, err := am.jwtService.ValidateToken(tokenString)
	if err != nil {
		panic(err)
	}

	ctx.Set(am.constants.Context.ID, uint(claims["sub"].(float64)))

	ctx.Next()
}
