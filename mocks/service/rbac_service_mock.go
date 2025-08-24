package mocks

import (
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type RBACServiceMock struct {
	mock.Mock
}

func NewRBACServiceMock() *RBACServiceMock {
	return &RBACServiceMock{}
}

func (r *RBACServiceMock) GetPermissions(request rbacdto.GetPermissionsListRequest) ([]rbacdto.PermissionResponse, int64, error) {
	args := r.Called(request)
	return args.Get(0).([]rbacdto.PermissionResponse), args.Get(1).(int64), args.Error(2)
}

func (r *RBACServiceMock) GetPermissionRoles(request rbacdto.GetPermissionRolesRequest) ([]rbacdto.RoleResponse, int64, error) {
	args := r.Called(request)
	return args.Get(0).([]rbacdto.RoleResponse), args.Get(1).(int64), args.Error(2)
}

func (r *RBACServiceMock) GetUserPermissions(user *entity.User) ([]rbacdto.PermissionResponse, error) {
	args := r.Called(user)
	return args.Get(0).([]rbacdto.PermissionResponse), args.Error(1)
}

func (r *RBACServiceMock) GetRoles(request rbacdto.GetRolesListRequest) ([]rbacdto.RoleResponse, int64, error) {
	args := r.Called(request)
	return args.Get(0).([]rbacdto.RoleResponse), args.Get(1).(int64), args.Error(2)
}

func (r *RBACServiceMock) CreateRole(newRoleRequest rbacdto.NewRoleRequest) error {
	args := r.Called(newRoleRequest)
	return args.Error(0)
}

func (r *RBACServiceMock) GetRoleDetails(roleID uint) (rbacdto.RoleResponse, error) {
	args := r.Called(roleID)
	return args.Get(0).(rbacdto.RoleResponse), args.Error(1)
}

func (r *RBACServiceMock) GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]*entity.User, int64, error) {
	args := r.Called(request)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (r *RBACServiceMock) GetUserRoles(user *entity.User) ([]rbacdto.RoleResponse, error) {
	args := r.Called(user)
	return args.Get(0).([]rbacdto.RoleResponse), args.Error(1)
}

func (r *RBACServiceMock) DeleteRole(roleID uint) error {
	args := r.Called(roleID)
	return args.Error(0)
}

func (r *RBACServiceMock) UpdateRole(newRoleRequest rbacdto.UpdateRoleRequest) error {
	args := r.Called(newRoleRequest)
	return args.Error(0)
}

func (r *RBACServiceMock) UpdateUserRoles(userRolesRequest rbacdto.UpdateUserRolesRequest) error {
	args := r.Called(userRolesRequest)
	return args.Error(0)
}

func (r *RBACServiceMock) UpdateStaffRoles(request rbacdto.UpdateStaffRolesRequest) error {
	args := r.Called(request)
	return args.Error(0)
}

func (r *RBACServiceMock) GetStaffRoles(staff *entity.CorporationStaff) ([]rbacdto.RoleResponse, error) {
	args := r.Called(staff)
	return args.Get(0).([]rbacdto.RoleResponse), args.Error(1)
}
