package userdto

import "mime/multipart"

type BasicRegisterRequest struct {
	FirstName string
	LastName  string
	Phone     string
	Password  string
}

type VerifyPhoneRequest struct {
	Phone string
	OTP   string
}

type LoginRequest struct {
	Phone    string
	Password string
}

type ForgotPasswordRequest struct {
	Phone string
}

type CompleteRegisterRequest struct {
	UserID       uint
	Email        string
	NationalCode string
	ProfilePic   *multipart.FileHeader
}

type ResetPasswordRequest struct {
	ID       uint
	Password string
}
