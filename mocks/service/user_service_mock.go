package mocks

import (
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (u *UserServiceMock) GetUserSortableColumns() []userdto.UserEnumResponse {
	args := u.Called()
	return args.Get(0).([]userdto.UserEnumResponse)
}

func (u *UserServiceMock) IsUserActive(userID uint) error {
	args := u.Called(userID)
	return args.Error(0)
}

func (u *UserServiceMock) GetUserByID(userID uint) (*entity.User, error) {
	args := u.Called(userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserServiceMock) GetUserCredential(userID uint) (userdto.CredentialResponse, error) {
	args := u.Called(userID)
	return args.Get(0).(userdto.CredentialResponse), args.Error(1)
}

func (u *UserServiceMock) GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	args := u.Called(permissionTypes)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserServiceMock) GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, int64, error) {
	args := u.Called(request)
	return args.Get(0).([]userdto.CredentialResponse), args.Get(1).(int64), args.Error(2)
}

func (u *UserServiceMock) BanUser(userID uint) error {
	args := u.Called(userID)
	return args.Error(0)
}

func (u *UserServiceMock) UnbanUser(userID uint) error {
	args := u.Called(userID)
	return args.Error(0)
}

func (u *UserServiceMock) Register(registerInfo userdto.BasicRegisterRequest) error {
	args := u.Called(registerInfo)
	return args.Error(0)
}

func (u *UserServiceMock) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error {
	args := u.Called(verifyInfo)
	return args.Error(0)
}

func (u *UserServiceMock) Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error) {
	args := u.Called(loginInfo)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}

func (u *UserServiceMock) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error {
	args := u.Called(forgotPasswordInfo)
	return args.Error(0)
}

func (u *UserServiceMock) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error) {
	args := u.Called(verifyInfo)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}

func (u *UserServiceMock) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error {
	args := u.Called(completeRegisterInfo)
	return args.Error(0)
}

func (u *UserServiceMock) VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error {
	args := u.Called(verifyOTPInfo)
	return args.Error(0)
}

func (u *UserServiceMock) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error {
	args := u.Called(resetPassInfo)
	return args.Error(0)
}

func (u *UserServiceMock) FindActiveUserByPhone(phone string) (*entity.User, error) {
	args := u.Called(phone)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserServiceMock) UpdateProfile(profileInfo userdto.UpdateProfileRequest) error {
	args := u.Called(profileInfo)
	return args.Error(0)
}

func (u *UserServiceMock) GetUserRoles(userID uint) ([]rbacdto.RoleResponse, error) {
	args := u.Called(userID)
	return args.Get(0).([]rbacdto.RoleResponse), args.Error(1)
}

func (u *UserServiceMock) UpdateUserRoles(request userdto.UpdateUserRolesRequest) error {
	args := u.Called(request)
	return args.Error(0)
}

func (u *UserServiceMock) GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]userdto.CredentialResponse, int64, error) {
	args := u.Called(request)
	return args.Get(0).([]userdto.CredentialResponse), args.Get(1).(int64), args.Error(2)
}

func (u *UserServiceMock) RefreshToken(refreshToken string) (userdto.UserInfoResponse, error) {
	args := u.Called(refreshToken)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}
