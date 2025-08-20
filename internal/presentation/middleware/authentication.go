package middleware

import (
	"strconv"
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
	constants             *bootstrap.Constants
	jwtService            usecase.JWTService
	userRepository        repository.UserRepository
	corporationRepository repository.CorporationRepository
	db                    database.Database
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	jwtService usecase.JWTService,
	userRepository repository.UserRepository,
	corporationRepository repository.CorporationRepository,
	db database.Database,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:             constants,
		jwtService:            jwtService,
		userRepository:        userRepository,
		corporationRepository: corporationRepository,
		db:                    db,
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

		if err := am.userRepository.FindUserRoles(am.db, user); err != nil {
			panic(err)
		}

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
		if err := am.userRepository.FindRolePermissions(am.db, &role); err != nil {
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

func (am *AuthMiddleware) RequireCorporationAccess(corporationIDField string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := ctx.Get(am.constants.Context.ID)

		corporationID, err := am.extractCorporationID(ctx, corporationIDField)
		if err != nil {
			panic(err)
		}

		if err := am.validateCorporationAccess(userID.(uint), uint(corporationID)); err != nil {
			panic(err)
		}

		ctx.Set(am.constants.Context.CorporationID, corporationID)
		ctx.Next()
	}
}

func (am *AuthMiddleware) RequireCorporationPermission(allowedPermissions []enum.PermissionType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := ctx.Get(am.constants.Context.ID)
		corporationID, _ := ctx.Get(am.constants.Context.CorporationID)

		staff, err := am.corporationRepository.FindCorporationStaff(am.db, userID.(uint), corporationID.(uint))
		if err != nil {
			panic(err)
		}

		if err = am.corporationRepository.FindStaffRoles(am.db, staff); err != nil {
			panic(err)
		}

		allowedPermissions = append(allowedPermissions, enum.PermissionAll, enum.CorporationPermissionAll)

		if !am.isAllowRole(allowedPermissions, staff.Roles) {
			err := exception.ForbiddenError{Resource: am.constants.Field.Page, Message: "access denied"}
			panic(err)
		}

		ctx.Next()
	}
}

func (am *AuthMiddleware) extractCorporationID(ctx *gin.Context, corporationIDField string) (uint, error) {
	var corporationID uint
	var found, invalid bool

	if corporationIDStr := ctx.Param(corporationIDField); corporationIDStr != "" {
		if id, err := am.parseUint(corporationIDStr); err == nil {
			corporationID = uint(id)
			found = true
		} else {
			invalid = true
		}
	}

	if !found {
		if corporationIDStr := ctx.Query(corporationIDField); corporationIDStr != "" {
			if id, err := am.parseUint(corporationIDStr); err == nil {
				corporationID = uint(id)
				found = true
			} else {
				invalid = true
			}
		}
	}

	if invalid {
		validationError := exception.ValidationErrors{}
		validationError.Add(corporationIDField, am.constants.Tag.InvalidNumber)
		return 0, validationError
	}

	if !found {
		validationError := exception.ValidationErrors{}
		validationError.Add(corporationIDField, am.constants.Tag.Required)
		return 0, validationError
	}

	if corporationID == 0 {
		validationError := exception.ValidationErrors{}
		validationError.Add(corporationIDField, am.constants.Tag.InvalidNumber)
		return 0, validationError
	}

	return corporationID, nil
}

func (am *AuthMiddleware) parseUint(str string) (uint, error) {
	number, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(number), nil
}

func (am *AuthMiddleware) validateCorporationAccess(userID, corporationID uint) error {
	staff, err := am.corporationRepository.FindCorporationStaff(am.db, userID, corporationID)
	if err != nil {
		return err
	}
	if staff == nil {
		return exception.NewNoPropertyAccessForbiddenError(am.constants.Field.Corporation)
	}
	return nil
}
