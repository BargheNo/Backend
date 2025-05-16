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

func (s *UserServiceMock) DoesUserExist(userID uint) {
	args := s.Called(userID)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) IsUserActive(userID uint) bool {
	args := s.Called(userID)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
	return args.Get(1).(bool)
}

func (s *UserServiceMock) GetUserByID(userID uint) *entity.User {
	args := s.Called(userID)
	return args.Get(0).(*entity.User)
}

func (s *UserServiceMock) GetUserCredential(userID uint) userdto.CredentialResponse {
	args := s.Called(userID)
	return args.Get(0).(userdto.CredentialResponse)
}

func (s *UserServiceMock) GetUsersByPermission(permissionTypes []enum.PermissionType) []*entity.User {
	args := s.Called(permissionTypes)
	return args.Get(0).([]*entity.User)
}

func (s *UserServiceMock) GetUsersByStatus(request userdto.GetUsersListRequest) []userdto.CredentialResponse {
	args := s.Called(request)
	return args.Get(0).([]userdto.CredentialResponse)
}

func (s *UserServiceMock) BanUser(userID uint) {
	banArgs := s.Called(userID)
	if banArgs.Get(0) != nil {
		panic(banArgs.Get(0))
	}
}

func (s *UserServiceMock) UnbanUser(userID uint) {
	unbanArgs := s.Called(userID)
	if unbanArgs.Get(0) != nil {
		panic(unbanArgs.Get(0))
	}
}

func (s *UserServiceMock) Register(registerInfo userdto.BasicRegisterRequest) {
	args := s.Called(registerInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) {
	args := s.Called(verifyInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse {
	args := s.Called(loginInfo)
	return args.Get(0).(userdto.UserInfoResponse)
}

func (s *UserServiceMock) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) {
	args := s.Called(forgotPasswordInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) userdto.UserInfoResponse {
	args := s.Called(verifyInfo)
	return args.Get(0).(userdto.UserInfoResponse)
}

func (s *UserServiceMock) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) {
	args := s.Called(completeRegisterInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) {
	args := s.Called(verifyOTPInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) {
	args := s.Called(resetPassInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (s *UserServiceMock) FindUserByPhone(phone string) userdto.UserResponse {
	args := s.Called(phone)
	return args.Get(0).(userdto.UserResponse)
}

func (s *UserServiceMock) UpdateProfile(profileInfo userdto.UpdateProfileRequest) {
	args := s.Called(profileInfo)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}

func (m *UserServiceMock) GetAllPermissions() []userdto.PermissionResponse {
	args := m.Called()
	return args.Get(0).([]userdto.PermissionResponse)
}

func (m *UserServiceMock) GetAllRoles() []userdto.RoleResponse {
	args := m.Called()
	return args.Get(0).([]userdto.RoleResponse)
}

func (m *UserServiceMock) CreateRole(newRoleRequest userdto.NewRoleRequest) {
	m.Called(newRoleRequest)
}

func (m *UserServiceMock) GetRoomDetails(roleID uint) userdto.RoleResponse {
	args := m.Called(roleID)
	return args.Get(0).(userdto.RoleResponse)
}

func (m *UserServiceMock) GetRoleOwners(roleID uint) []userdto.CredentialResponse {
	args := m.Called(roleID)
	return args.Get(0).([]userdto.CredentialResponse)
}

func (m *UserServiceMock) GetUserRoles(userID uint) []userdto.RoleResponse {
	args := m.Called(userID)
	return args.Get(0).([]userdto.RoleResponse)
}

func (m *UserServiceMock) DeleteRole(roleID uint) {
	m.Called(roleID)
}

func (m *UserServiceMock) UpdateRole(newRoleRequest userdto.UpdateRoleRequest) {
	m.Called(newRoleRequest)
}

func (m *UserServiceMock) UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) {
	m.Called(userRolesRequest)
}
