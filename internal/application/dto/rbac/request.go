package rbacdto

import "github.com/BargheNo/Backend/internal/domain/entity"

type GetRolesListRequest struct {
	IsStaff bool
	Query   string
	Offset  int
	Limit   int
}

type GetPermissionsListRequest struct {
	IsStaff bool
	Offset  int
	Limit   int
}

type NewRoleRequest struct {
	Name          string
	IsStaff       bool
	PermissionIDs []uint
}

type GetRoleOwnersRequest struct {
	RoleID uint
	Query  string
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type UpdateRoleRequest struct {
	RoleID        uint
	Name          *string
	PermissionIDs []uint
}

type UpdateUserRolesRequest struct {
	User    *entity.User
	RoleIDs []uint
}

type UpdateStaffRolesRequest struct {
	Staff   *entity.CorporationStaff
	RoleIDs []uint
}

type GetPermissionRolesRequest struct {
	PermissionID uint
	Offset       int
	Limit        int
	SortBy       uint
	Asc          bool
}
