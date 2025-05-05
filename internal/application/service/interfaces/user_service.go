package service

import (
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type UserService interface {
	DoesUserExist(userID uint)
	IsUserActive(userID uint) bool
	GetUserByID(userID uint) *entity.User
	GetUserCredential(userID uint) userdto.CredentialResponse
	GetUsersByPermission(permissionTypes []enum.PermissionType) []*entity.User
	GetUsersByStatus(request userdto.GetUsersListRequest) []userdto.CredentialResponse
	BanUser(userID uint)
	UnbanUser(userID uint)
	Register(registerInfo userdto.BasicRegisterRequest)
	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest)
	Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse
	ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest)
	VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) userdto.UserInfoResponse
	CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest)
	VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest)
	ResetPassword(resetPassInfo userdto.ResetPasswordRequest)
	FindUserByPhone(phone string) userdto.UserResponse
	UpdateProfile(profileInfo userdto.UpdateProfileRequest)
	GetAllPermissions() []userdto.PermissionResponse
	GetAllRoles() []userdto.RoleResponse
	CreateRole(newRoleRequest userdto.NewRoleRequest)
	GetRoomDetails(roleID uint) userdto.RoleResponse
	GetRoleOwners(roleID uint) []userdto.CredentialResponse
	GetUserRoles(userID uint) []userdto.RoleResponse
	DeleteRole(roleID uint)
	UpdateRole(newRoleRequest userdto.UpdateRoleRequest)
	UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest)
}
