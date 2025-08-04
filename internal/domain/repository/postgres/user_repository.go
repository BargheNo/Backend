package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUsers(db database.Database) ([]*entity.User, error)
	FindUserByID(db database.Database, id uint) (*entity.User, error)
	FindUserByPhone(db database.Database, phone string) (*entity.User, error)
	FindUserByEmail(db database.Database, email string) (*entity.User, error)
	CreateUser(db database.Database, user *entity.User) error
	DeleteUserByPhone(db database.Database, phone string) error
	UpdateUser(db database.Database, user *entity.User) error
	FindUserRoles(db database.Database, user *entity.User) error
	FindRoleByName(db database.Database, name string) (*entity.Role, error)
	FindRolePermissions(db database.Database, role *entity.Role) error
	FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, error)
	RoleHasPermission(db database.Database, roleID uint, permissionID uint) bool
	UserHasRole(db database.Database, userID uint, roleID uint) bool
	CreateRole(db database.Database, role *entity.Role) error
	CreatePermission(db database.Database, permission *entity.Permission) error
	AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error
	AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error
	FindAllPermissions(db database.Database) ([]*entity.Permission, error)
	FindAllRoles(db database.Database) ([]*entity.Role, error)
	FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error)
	FindRolesByPermission(db database.Database, permissionID uint, opts ...QueryModifier) ([]*entity.Role, error)
	FindRoleByID(db database.Database, roleID uint) (*entity.Role, error)
	FindUsersByRoleID(db database.Database, roleID uint) ([]*entity.User, error)
	FindUserByStatus(db database.Database, status []enum.UserStatus, opts ...QueryModifier) ([]*entity.User, error)
	FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) ([]*entity.User, error)
	DeleteRole(db database.Database, roleID uint) error
	UpdateRole(db database.Database, role *entity.Role) error
	ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error
	ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error
}
