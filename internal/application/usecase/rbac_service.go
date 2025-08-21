package usecase

import (
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type RBACService interface {
	GetPermissions(request rbacdto.GetPermissionsListRequest) ([]rbacdto.PermissionResponse, int64, error)
	GetPermissionRoles(request rbacdto.GetPermissionRolesRequest) ([]rbacdto.RoleResponse, int64, error)
	GetUserPermissions(user *entity.User) ([]rbacdto.PermissionResponse, error)
	GetRoles(request rbacdto.GetRolesListRequest) ([]rbacdto.RoleResponse, int64, error)
	CreateRole(newRoleRequest rbacdto.NewRoleRequest) error
	GetRoleDetails(roleID uint) (rbacdto.RoleResponse, error)
	GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]*entity.User, int64, error)
	GetUserRoles(user *entity.User) ([]rbacdto.RoleResponse, error)
	DeleteRole(roleID uint) error
	UpdateRole(newRoleRequest rbacdto.UpdateRoleRequest) error
	UpdateUserRoles(userRolesRequest rbacdto.UpdateUserRolesRequest) error
	UpdateStaffRoles(request rbacdto.UpdateStaffRolesRequest) error
	GetStaffRoles(staff *entity.CorporationStaff) ([]rbacdto.RoleResponse, error)
}
