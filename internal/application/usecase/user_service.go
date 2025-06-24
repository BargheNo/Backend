package usecase

import (
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type UserService interface {
	IsUserActive(userID uint) error
	GetUserByID(userID uint) (*entity.User, error)
	GetUserCredential(userID uint) (userdto.CredentialResponse, error)
	GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error)
	GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, error)
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
	GetAllPermissions() ([]userdto.PermissionResponse, error)
	GetAllRoles() ([]userdto.RoleResponse, error)
	CreateRole(newRoleRequest userdto.NewRoleRequest) error
	GetRoleDetails(roleID uint) (userdto.RoleResponse, error)
	GetRoleOwners(roleID uint) ([]userdto.CredentialResponse, error)
	GetUserRoles(userID uint) ([]userdto.RoleResponse, error)
	DeleteRole(roleID uint) error
	UpdateRole(newRoleRequest userdto.UpdateRoleRequest) error
	UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) error
}
