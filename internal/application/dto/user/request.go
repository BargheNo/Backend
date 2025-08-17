package userdto

import "mime/multipart"

type BasicRegisterRequest struct {
	FirstName string
	LastName  string
	Phone     string
	Password  string
	Recaptcha string
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
	Phone     string
	Password  string
	Recaptcha string
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
	IsStaff       bool
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
	Query  string
	Status uint
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type GetPermissionRolesRequest struct {
	PermissionID uint
	Offset       int
	Limit        int
	SortBy       uint
	Asc          bool
}

type GetPermissionsListRequest struct {
	Offset int
	Limit  int
}

type GetRoleOwnersRequest struct {
	RoleID uint
	Offset int
	Limit  int
}

type GetRolesListRequest struct {
	Query  string
	Offset int
	Limit  int
}

type SearchUsersRequest struct {
	Query  string
	Status uint
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}
