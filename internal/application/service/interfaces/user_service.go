package service

import userdto "github.com/BargheNo/Backend/internal/application/dto/user"

type UserService interface {
	Register(registerInfo userdto.BasicRegisterRequest)
	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest)
	Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse
	ForgotPassword(forgotPasswordInfo userdto.ForgotPassword)
}
