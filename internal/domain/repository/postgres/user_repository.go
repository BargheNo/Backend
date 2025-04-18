package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUserByID(db database.Database, id uint) (*entity.User, bool)
	FindUserByPhone(db database.Database, phone string) (*entity.User, bool)
	FindUserByEmail(db database.Database, email string) (*entity.User, bool)
	CreateUser(db database.Database, user *entity.User) error
	DeleteUserByPhone(db database.Database, phone string) error
	UpdateUser(db database.Database, user *entity.User) error
	FindUserRoles(db database.Database, user *entity.User) error
	FindRoleByName(db database.Database, name string) (*entity.Role, bool)
	FindRolePermissions(db database.Database, role *entity.Role) error
	FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, bool)
	RoleHasPermission(db database.Database, roleID uint, permissionID uint) bool
	UserHasRole(db database.Database, userID uint, roleID uint) bool
	CreateRole(db database.Database, role *entity.Role) error
	CreatePermission(db database.Database, permission *entity.Permission) error
	AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error
	AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error
}
