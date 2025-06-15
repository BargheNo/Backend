package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

func (u *UserRepositoryMock) FindUsers(db database.Database) ([]*entity.User, error) {
	args := u.Called(db)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByID(db database.Database, id uint) (*entity.User, error) {
	args := u.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByPhone(db database.Database, phone string) (*entity.User, error) {
	args := u.Called(db, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByEmail(db database.Database, email string) (*entity.User, error) {
	args := u.Called(db, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) CreateUser(db database.Database, user *entity.User) error {
	args := u.Called(db, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) DeleteUserByPhone(db database.Database, phone string) error {
	args := u.Called(db, phone)
	return args.Error(0)
}

func (u *UserRepositoryMock) UpdateUser(db database.Database, user *entity.User) error {
	args := u.Called(db, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) FindUserRoles(db database.Database, user *entity.User) error {
	args := u.Called(db, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	args := u.Called(db, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) FindRolePermissions(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(0)
}

func (u *UserRepositoryMock) FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, error) {
	args := u.Called(db, permissionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) RoleHasPermission(db database.Database, roleID uint, permissionID uint) bool {
	args := u.Called(db, roleID, permissionID)
	return args.Bool(0)
}

func (u *UserRepositoryMock) UserHasRole(db database.Database, userID uint, roleID uint) bool {
	args := u.Called(db, userID, roleID)
	return args.Bool(0)
}

func (u *UserRepositoryMock) CreateRole(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(0)
}

func (u *UserRepositoryMock) CreatePermission(db database.Database, permission *entity.Permission) error {
	args := u.Called(db, permission)
	return args.Error(0)
}

func (u *UserRepositoryMock) AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error {
	args := u.Called(db, role, permission)
	return args.Error(0)
}

func (u *UserRepositoryMock) AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error {
	args := u.Called(db, user, role)
	return args.Error(0)
}

func (u *UserRepositoryMock) FindAllPermissions(db database.Database) ([]*entity.Permission, error) {
	args := u.Called(db)
	return args.Get(0).([]*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) FindAllRoles(db database.Database) ([]*entity.Role, error) {
	args := u.Called(db)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error) {
	args := u.Called(db, permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) FindRoleByID(db database.Database, roleID uint) (*entity.Role, error) {
	args := u.Called(db, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) FindUsersByRoleID(db database.Database, roleID uint) ([]*entity.User, error) {
	args := u.Called(db, roleID)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByStatus(db database.Database, status []enum.UserStatus, opts ...repository.QueryModifier) ([]*entity.User, error) {
	args := u.Called(db, status, opts)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	args := u.Called(db, permissionTypes)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) DeleteRole(db database.Database, roleID uint) error {
	args := u.Called(db, roleID)
	return args.Error(0)
}

func (u *UserRepositoryMock) UpdateRole(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(0)
}

func (u *UserRepositoryMock) ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error {
	args := u.Called(db, role, permissions)
	return args.Error(0)
}

func (u *UserRepositoryMock) ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error {
	args := u.Called(db, user, roles)
	return args.Error(0)
}
