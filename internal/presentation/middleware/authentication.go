package middleware

import (
	"strings"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	jwtService     usecase.JWTService
	userRepository repository.UserRepository
	db             database.Database
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	jwtService usecase.JWTService,
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

func (am *AuthMiddleware) RequiredWithPermission(allowedPermissions []enum.PermissionType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, exist := ctx.Get(am.constants.Context.ID)
		if !exist {
			unauthorizedError := exception.NewUnauthorizedError("", nil)
			panic(unauthorizedError)
		}
		user, _ := am.userRepository.FindUserByID(am.db, id.(uint))
		am.userRepository.FindUserRoles(am.db, user)
		allowedPermissions = append(allowedPermissions, enum.PermissionAll)
		if !am.isAllowRole(allowedPermissions, user.Roles) {
			err := exception.ForbiddenError{Resource: am.constants.Field.Page, Message: "access denied"}
			panic(err)
		}
		ctx.Next()
	}
}

func (am *AuthMiddleware) isAllowRole(allowedPermissions []enum.PermissionType, roles []entity.Role) bool {
	allowedPermissionMap := make(map[enum.PermissionType]bool)
	for _, permission := range allowedPermissions {
		allowedPermissionMap[permission] = true
	}
	for _, role := range roles {
		err := am.userRepository.FindRolePermissions(am.db, &role)
		if err != nil {
			panic(err)
		}
		for _, permission := range role.Permissions {
			if allowedPermissionMap[permission.Type] {
				return true
			}
		}
	}
	return false
}
