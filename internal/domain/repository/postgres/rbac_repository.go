package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type RBACRepository interface {
	FindAllPermissions(db database.Database, userType enum.UserType, options *QueryOptions) ([]*entity.Permission, error)
	CountAllPermissions(db database.Database, userType enum.UserType) (int64, error)
	FindRolePermissions(db database.Database, role *entity.Role) error
	FindAllRoles(db database.Database, userType enum.UserType, options *QueryOptions) ([]*entity.Role, error)
	CountAllRoles(db database.Database, userType enum.UserType) (int64, error)
	FindRolesByQuery(db database.Database, userType enum.UserType, query string, options *QueryOptions) ([]*entity.Role, error)
	CountRolesByQuery(db database.Database, userType enum.UserType, query string) (int64, error)
	FindRoleByName(db database.Database, name string) (*entity.Role, error)
	CreateRole(db database.Database, role *entity.Role) error
	FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error)
	FindRolesByPermission(db database.Database, permissionID uint, options *QueryOptions) ([]*entity.Role, error)
	CountRolesByPermission(db database.Database, permissionID uint) (int64, error)
	AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error
	FindRoleByID(db database.Database, roleID uint) (*entity.Role, error)
	FindUsersByRoleID(db database.Database, roleID uint, options *QueryOptions) ([]*entity.User, error)
	CountUsersByRoleID(db database.Database, roleID uint) (int64, error)
	FindUsersByRoleIDAndQuery(db database.Database, roleID uint, query string, options *QueryOptions) ([]*entity.User, error)
	CountUsersByRoleIDAndQuery(db database.Database, roleID uint, query string) (int64, error)
	FindUserRoles(db database.Database, user *entity.User) error
	DeleteRole(db database.Database, roleID uint) error
	UpdateRole(db database.Database, role *entity.Role) error
	ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []*entity.Permission) error
	ReplaceUserRoles(db database.Database, user *entity.User, roles []*entity.Role) error
	ReplaceStaffRoles(db database.Database, staff *entity.CorporationStaff, roles []*entity.Role) error
	FindStaffRoles(db database.Database, staff *entity.CorporationStaff) error
}
