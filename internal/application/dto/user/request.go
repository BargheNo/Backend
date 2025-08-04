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

type VerifyEmailRequest struct {
	UserID uint
	Email  string
	OTP    string
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
	TemplateFile string
	EmailSubject string
}

type ResetPasswordRequest struct {
	UserID   uint
	Password string
}

type UpdateProfileRequest struct {
	UserID       uint
	FirstName    *string
	LastName     *string
	Email        *string
	NationalCode *string
	ProfilePic   *multipart.FileHeader
	TemplateFile string
	EmailSubject string
}

type NewRoleRequest struct {
	Name          string
	PermissionIDs []uint
}

type UpdateRoleRequest struct {
	RoleID        uint
	Name          *string
	PermissionIDs []uint
}

type UpdateUserRolesRequest struct {
	UserID  uint
	RoleIDs []uint
}

type GetUsersListRequest struct {
	Statuses []uint
	Offset   int
	Limit    int
}

type GetPermissionRolesRequest struct {
	PermissionID uint
	Offset       int
	Limit        int
}
