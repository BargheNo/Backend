package middleware

import (
	"strings"

	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants  *bootstrap.Constants
	jwtService service.JWTService
	db         database.Database
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:  constants,
		jwtService: jwtService,
		db:         db,
	}
}

func (am *AuthMiddleware) AuthRequired(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		unauthorizedError:= exception.NewUnauthorizedError("Missing Authorization header", nil)
		panic(unauthorizedError)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError := exception.NewUnauthorizedError("Invalid Authorization format", nil)
		panic(unauthorizedError)
	}

	tokenString := parts[1]
	if tokenString == "" {
		unauthorizedError := exception.NewUnauthorizedError("Missing token", nil)
		panic(unauthorizedError)
	}
	claims, err := am.jwtService.ValidateToken(tokenString)
	if err != nil {
		unauthorizedError := exception.NewUnauthorizedError("Invalid token", err)
		panic(unauthorizedError)
	}

	c.Set(am.constants.Context.ID, uint(claims["sub"].(float64)))

	c.Next()
}
