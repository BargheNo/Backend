package mocks

import (
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (s *UserServiceMock) GetUserCredential(userID uint) userdto.CredentialResponse {
	args := s.Called(userID)
	return args.Get(0).(userdto.CredentialResponse)
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
