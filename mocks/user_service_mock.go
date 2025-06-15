package mocks

import (
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

func (s *UserServiceMock) DoesUserExist(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) IsUserActive(userID uint) (bool, error) {
	args := s.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (s *UserServiceMock) GetUserByID(userID uint) (*entity.User, error) {
	args := s.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (s *UserServiceMock) GetUserCredential(userID uint) (userdto.CredentialResponse, error) {
	args := s.Called(userID)
	return args.Get(0).(userdto.CredentialResponse), args.Error(1)
}

func (s *UserServiceMock) GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	args := s.Called(permissionTypes)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (s *UserServiceMock) GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, error) {
	args := s.Called(request)
	return args.Get(0).([]userdto.CredentialResponse), args.Error(1)
}

func (s *UserServiceMock) BanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) UnbanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) Register(registerInfo userdto.BasicRegisterRequest) error {
	args := s.Called(registerInfo)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error {
	args := s.Called(verifyInfo)
	return args.Error(0)
}

func (s *UserServiceMock) Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error) {
	args := s.Called(loginInfo)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}

func (s *UserServiceMock) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error {
	args := s.Called(forgotPasswordInfo)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error) {
	args := s.Called(verifyInfo)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}

func (s *UserServiceMock) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error {
	args := s.Called(completeRegisterInfo)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error {
	args := s.Called(verifyOTPInfo)
	return args.Error(0)
}

func (s *UserServiceMock) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error {
	args := s.Called(resetPassInfo)
	return args.Error(0)
}

func (s *UserServiceMock) FindUserByPhone(phone string) (userdto.UserResponse, error) {
	args := s.Called(phone)
	return args.Get(0).(userdto.UserResponse), args.Error(1)
}

func (s *UserServiceMock) UpdateProfile(profileInfo userdto.UpdateProfileRequest) error {
	args := s.Called(profileInfo)
	return args.Error(0)
}

func (s *UserServiceMock) GetAllPermissions() ([]userdto.PermissionResponse, error) {
	args := s.Called()
	return args.Get(0).([]userdto.PermissionResponse), args.Error(1)
}

func (s *UserServiceMock) GetAllRoles() ([]userdto.RoleResponse, error) {
	args := s.Called()
	return args.Get(0).([]userdto.RoleResponse), args.Error(1)
}

func (s *UserServiceMock) CreateRole(newRoleRequest userdto.NewRoleRequest) error {
	args := s.Called(newRoleRequest)
	return args.Error(0)
}

func (s *UserServiceMock) GetRoomDetails(roleID uint) (userdto.RoleResponse, error) {
	args := s.Called(roleID)
	return args.Get(0).(userdto.RoleResponse), args.Error(1)
}

func (s *UserServiceMock) GetRoleOwners(roleID uint) ([]userdto.CredentialResponse, error) {
	args := s.Called(roleID)
	return args.Get(0).([]userdto.CredentialResponse), args.Error(1)
}

func (s *UserServiceMock) GetUserRoles(userID uint) ([]userdto.RoleResponse, error) {
	args := s.Called(userID)
	return args.Get(0).([]userdto.RoleResponse), args.Error(1)
}

func (s *UserServiceMock) DeleteRole(roleID uint) error {
	args := s.Called(roleID)
	return args.Error(0)
}

func (s *UserServiceMock) UpdateRole(newRoleRequest userdto.UpdateRoleRequest) error {
	args := s.Called(newRoleRequest)
	return args.Error(0)
}

func (s *UserServiceMock) UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) error {
	args := s.Called(userRolesRequest)
	return args.Error(0)
}
