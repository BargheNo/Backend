package service

import userdto "github.com/BargheNo/Backend/internal/application/dto/user"

type UserService interface {
	GetUserCredential(userID uint) userdto.CredentialResponse
	Register(registerInfo userdto.BasicRegisterRequest)
	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest)
	Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse
	ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest)
	VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) userdto.UserInfoResponse
	ResetPassword(resetPassInfo userdto.ResetPasswordRequest)
	FindUserByPhone(phone string) userdto.CredentialResponse
}
