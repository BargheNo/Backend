package usecase

import (
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type UserService interface {
	GetUserSortableColumns() []userdto.UserEnumResponse
	IsUserActive(userID uint) error
	GetUserByID(userID uint) (*entity.User, error)
	GetUserCredential(userID uint) (userdto.CredentialResponse, error)
	GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error)
	GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, int64, error)
	BanUser(userID uint) error
	UnbanUser(userID uint) error
	Register(registerInfo userdto.BasicRegisterRequest) error
	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error
	Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error)
	ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error
	VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error)
	CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error
	VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error
	ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error
	FindActiveUserByPhone(phone string) (*entity.User, error)
	UpdateProfile(profileInfo userdto.UpdateProfileRequest) error
	GetUserRoles(userID uint) ([]rbacdto.RoleResponse, error)
	UpdateUserRoles(request userdto.UpdateUserRolesRequest) error
	GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]userdto.CredentialResponse, int64, error)
	RefreshToken(refreshToken string) (userdto.UserInfoResponse, error)
}
