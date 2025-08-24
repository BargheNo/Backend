package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByStatus(db database.Database, statuses []enum.UserStatus, options *postgres.QueryOptions) ([]*entity.User, error) {
	args := u.Called(db, statuses, options)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) CountUserByStatus(db database.Database, statuses []enum.UserStatus) (int64, error) {
	args := u.Called(db, statuses)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByEmail(db database.Database, email string) (*entity.User, error) {
	args := u.Called(db, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByPhone(db database.Database, phone string) (*entity.User, error) {
	args := u.Called(db, phone)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	args := u.Called(db, name)
	return args.Get(0).(*entity.Role), args.Error(1)
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
	return args.Error(1)
}

func (u *UserRepositoryMock) FindRolePermissions(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(1)
}

func (u *UserRepositoryMock) FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, error) {
	args := u.Called(db, permissionType)
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) RoleHasPermission(db database.Database, roleID uint, permissionID uint) bool {
	args := u.Called(db, roleID, permissionID)
	return args.Get(0).(bool)
}

func (u *UserRepositoryMock) UserHasRole(db database.Database, userID uint, roleID uint) bool {
	args := u.Called(db, userID, roleID)
	return args.Get(0).(bool)
}

func (u *UserRepositoryMock) CreateRole(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(1)
}

func (u *UserRepositoryMock) CreatePermission(db database.Database, permission *entity.Permission) error {
	args := u.Called(db, permission)
	return args.Error(1)
}

func (u *UserRepositoryMock) AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error {
	args := u.Called(db, role, permission)
	return args.Error(1)
}

func (u *UserRepositoryMock) AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error {
	args := u.Called(db, user, role)
	return args.Error(1)
}

func (u *UserRepositoryMock) FindAllPermissions(db database.Database, userType enum.UserType, options *postgres.QueryOptions) ([]*entity.Permission, error) {
	args := u.Called(db, userType, options)
	return args.Get(0).([]*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) CountAllPermissions(db database.Database, userType enum.UserType) (int64, error) {
	args := u.Called(db, userType)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) FindAllRoles(db database.Database, userType enum.UserType, options *postgres.QueryOptions) ([]*entity.Role, error) {
	args := u.Called(db, userType, options)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) FindRolesByQuery(db database.Database, userType enum.UserType, query string, options *postgres.QueryOptions) ([]*entity.Role, error) {
	args := u.Called(db, userType, query, options)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) CountRolesByQuery(db database.Database, userType enum.UserType, query string) (int64, error) {
	args := u.Called(db, userType, query)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) CountAllRoles(db database.Database, userType enum.UserType) (int64, error) {
	args := u.Called(db, userType)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error) {
	args := u.Called(db, permissionID)
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (u *UserRepositoryMock) FindRoleByID(db database.Database, roleID uint) (*entity.Role, error) {
	args := u.Called(db, roleID)
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) FindUsersByRoleID(db database.Database, roleID uint, options *postgres.QueryOptions) ([]*entity.User, error) {
	args := u.Called(db, roleID, options)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) CountUsersByRoleID(db database.Database, roleID uint) (int64, error) {
	args := u.Called(db, roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	args := u.Called(db, permissionTypes)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) DeleteRole(db database.Database, roleID uint) error {
	args := u.Called(db, roleID)
	return args.Error(1)
}

func (u *UserRepositoryMock) UpdateRole(db database.Database, role *entity.Role) error {
	args := u.Called(db, role)
	return args.Error(1)
}

func (u *UserRepositoryMock) ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error {
	args := u.Called(db, role, permissions)
	return args.Error(1)
}

func (u *UserRepositoryMock) ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error {
	args := u.Called(db, user, roles)
	return args.Error(1)
}

func (u *UserRepositoryMock) FindRolesByPermission(db database.Database, permissionID uint, options *postgres.QueryOptions) ([]*entity.Role, error) {
	args := u.Called(db, permissionID, options)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (u *UserRepositoryMock) CountRolesByPermission(db database.Database, permissionID uint) (int64, error) {
	args := u.Called(db, permissionID)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepositoryMock) FindUsersByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.User, error) {
	args := u.Called(db, query, options)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) CountUsersByQuery(db database.Database, query string) (int64, error) {
	args := u.Called(db, query)
	return args.Get(0).(int64), args.Error(1)
}
