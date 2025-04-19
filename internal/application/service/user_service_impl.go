package serviceimpl

import (
	"context"
	"regexp"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	loggerimpl "github.com/BargheNo/Backend/internal/application/adapter/logger"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/logger"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	cacherepository "github.com/BargheNo/Backend/internal/domain/repository/redis"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants           *bootstrap.Constants
	otpService          service.OTPService
	jwtService          service.JWTService
	smsService          service.SMSService
	emailService        service.EmailService
	s3Storage           s3.S3Storage
	userRepository      repository.UserRepository
	userCacheRepository cacherepository.UserCacheRepository
	db                  database.Database
}

type UserServiceDeps struct {
	Constants           *bootstrap.Constants
	OTPService          service.OTPService
	JWTService          service.JWTService
	SMSService          service.SMSService
	EmailService        service.EmailService
	S3Storage           s3.S3Storage
	UserRepository      repository.UserRepository
	UserCacheRepository cacherepository.UserCacheRepository
	DB                  database.Database
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:           deps.Constants,
		otpService:          deps.OTPService,
		jwtService:          deps.JWTService,
		smsService:          deps.SMSService,
		emailService:        deps.EmailService,
		s3Storage:           deps.S3Storage,
		userRepository:      deps.UserRepository,
		userCacheRepository: deps.UserCacheRepository,
		db:                  deps.DB,
	}
}

func (userService *UserService) validatePasswordTests(errors *[]string, test string, password string, tag string) {
	matched, _ := regexp.MatchString(test, password)
	if !matched {
		*errors = append(*errors, tag)
	}
}

func (userService *UserService) passwordValidation(password string) error {
	var errors exception.ValidationErrors
	var errorTags []string

	userService.validatePasswordTests(&errorTags, ".{8,}", password, userService.constants.Tag.MinimumLength)
	userService.validatePasswordTests(&errorTags, "[a-z]", password, userService.constants.Tag.ContainsLowercase)
	userService.validatePasswordTests(&errorTags, "[A-Z]", password, userService.constants.Tag.ContainsUppercase)
	userService.validatePasswordTests(&errorTags, "[0-9]", password, userService.constants.Tag.ContainsNumber)
	userService.validatePasswordTests(&errorTags, "[^\\d\\w]", password, userService.constants.Tag.ContainsSpecialChar)

	for _, tag := range errorTags {
		errors.Add(userService.constants.Field.Password, tag)
	}
	if len(errorTags) > 0 {
		return errors
	}
	return nil
}

func (userService *UserService) validateDuplicateEmail(email string) error {
	var conflictErrors exception.ConflictErrors
	redisKey := userService.constants.RedisKey.GenerateOTPKey(email)
	_, exist := userService.userCacheRepository.Get(context.Background(), redisKey)
	if exist {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, userExist := userService.userRepository.FindUserByEmail(userService.db, email)
	if userExist && user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) validateDuplicatePhone(phone string) error {
	var conflictErrors exception.ConflictErrors
	redisKey := userService.constants.RedisKey.GenerateOTPKey(phone)
	_, exist := userService.userCacheRepository.Get(context.Background(), redisKey)
	if exist {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, userExist := userService.userRepository.FindUserByPhone(userService.db, phone)
	if userExist && user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) enterNewEmail(firstName, lastName, email, emailSubject, templateFile string) {
	err := userService.validateDuplicateEmail(email)
	if err != nil {
		panic(err)
	}

	otp, expiryMinute := userService.otpService.GenerateOTP()
	redisKey := userService.constants.RedisKey.GenerateOTPKey(email)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		panic(err)
	}

	data := struct {
		FirstName    string
		LastName     string
		OTP          string
		ExpiryMinute int
		Year         int
	}{
		FirstName:    firstName,
		LastName:     lastName,
		OTP:          otp,
		ExpiryMinute: expiryMinute,
		Year:         time.Now().Year(),
	}
	userService.emailService.SendEmail(email, emailSubject, templateFile, data)
}

func (userService *UserService) DoesUserExist(userID uint) {
	_, userExist := userService.userRepository.FindUserByID(userService.db, userID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
}

func (userService *UserService) IsUserActive(userID uint) bool {
	user, userExist := userService.userRepository.FindUserByID(userService.db, userID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	isActive := enum.UserStatusActive == user.Status
	return isActive
}

func (userService *UserService) GetUserCredential(userID uint) userdto.CredentialResponse {
	user, userExist := userService.userRepository.FindUserByID(userService.db, userID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	profilePic := ""
	if user.ProfilePicPath != "" {
		profilePic = userService.s3Storage.GetPresignedURL(enum.ProfilePic, user.ProfilePicPath, 8*time.Hour)
	}
	return userdto.CredentialResponse{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		Email:      user.Email,
		NationalID: user.NationalCode,
		ProfilePic: profilePic,
		Status:     user.Status.String(),
	}
}

func (userService *UserService) Register(registerInfo userdto.BasicRegisterRequest) {
	err := userService.validateDuplicatePhone(registerInfo.Phone)
	if err != nil {
		panic(err)
	}

	err = userService.passwordValidation(registerInfo.Password)
	if err != nil {
		panic(err)
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), 14)
	if err != nil {
		panic(err)
	}

	err = userService.userRepository.DeleteUserByPhone(userService.db, registerInfo.Phone)
	if err != nil {
		panic(err)
	}

	user := &entity.User{
		FirstName:     registerInfo.FirstName,
		LastName:      registerInfo.LastName,
		Phone:         registerInfo.Phone,
		Password:      string(hashesPasswordBytes),
		PhoneVerified: false,
		EmailVerified: false,
		Status:        enum.UserStatusActive,
	}
	err = userService.userRepository.CreateUser(userService.db, user)
	if err != nil {
		panic(err)
	}

	otp, expiryMinute := userService.otpService.GenerateOTP()
	redisKey := userService.constants.RedisKey.GenerateOTPKey(registerInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		panic(err)
	}
	// userService.smsService.SendOTP(registerInfo.Phone, otp)
}

func (userService *UserService) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, verifyInfo.Phone)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err := userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		panic(err)
	}
	user.PhoneVerified = true
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		panic(err)
	}
}

func (userService *UserService) FindUserPermissions(user *entity.User) []userdto.PermissionResponse {
	var permissions []userdto.PermissionResponse
	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		panic(err)
	}
	for _, role := range user.Roles {
		rolePermissions := userService.getRolePermissions(&role)
		permissions = append(permissions, rolePermissions...)
	}
	return permissions
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, loginInfo.Phone)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		authError := exception.NewInvalidCredentialsError("phone and password not match", nil)
		panic(authError)
	}
	accessToken, refreshToken := userService.jwtService.GenerateToken(user.ID)
	permissions := userService.FindUserPermissions(user)
	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Permissions:  permissions,
	}
}

func (userService *UserService) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, forgotPasswordInfo.Phone)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}
	otp, expiryMinute := userService.otpService.GenerateOTP()
	redisKey := userService.constants.RedisKey.GenerateOTPKey(forgotPasswordInfo.Phone)
	err := userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		panic(err)
	}
	// userService.smsService.SendOTP(registerInfo.Phone, otp)
}

func (userService *UserService) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) userdto.UserInfoResponse {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, verifyInfo.Phone)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err := userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		panic(err)
	}

	accessToken, refreshToken := userService.jwtService.GenerateToken(user.ID)
	permissions := userService.FindUserPermissions(user)
	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Permissions:  permissions,
	}
}

func (userService *UserService) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) {
	user, userExist := userService.userRepository.FindUserByID(userService.db, completeRegisterInfo.UserID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if completeRegisterInfo.Email != "" {
		userService.enterNewEmail(user.FirstName, user.LastName, completeRegisterInfo.Email, completeRegisterInfo.EmailSubject, completeRegisterInfo.TemplateFile)
	}
	user.Email = completeRegisterInfo.Email
	user.EmailVerified = false
	user.NationalCode = completeRegisterInfo.NationalCode
	if completeRegisterInfo.ProfilePic != nil {
		profilePicPath := userService.constants.S3BucketPath.GetUserProfilePath(completeRegisterInfo.UserID, completeRegisterInfo.ProfilePic.Filename)
		userService.s3Storage.UploadObject(enum.ProfilePic, profilePicPath, completeRegisterInfo.ProfilePic)
		user.ProfilePicPath = profilePicPath
	}
	err := userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		panic(err)
	}
}

func (userService *UserService) VerifyEmail(verifyInfo userdto.VerifyEmailRequest) {
	var conflictErrors exception.ConflictErrors
	user, userExist := userService.userRepository.FindUserByID(userService.db, verifyInfo.UserID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}
	if user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Email)
	err := userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		panic(err)
	}
	user.EmailVerified = true
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		panic(err)
	}
}

func (userService *UserService) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) {
	user, userExist := userService.userRepository.FindUserByID(userService.db, resetPassInfo.ID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}

	if err := userService.passwordValidation(resetPassInfo.Password); err != nil {
		panic(err)
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(resetPassInfo.Password), 14)
	if err != nil {
		panic(err)
	}

	user.Password = string(hashesPasswordBytes)
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		panic(err)
	}
}

func (userService *UserService) FindUserByPhone(phone string) userdto.UserResponse {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, phone)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		panic(conflictErrors)
	}
	return userdto.UserResponse{
		ID: user.ID,
	}
}

func (userService *UserService) UpdateProfile(profileInfo userdto.UpdateProfileRequest) {
	user, userExist := userService.userRepository.FindUserByID(userService.db, profileInfo.UserID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}

	if profileInfo.FirstName != nil {
		user.FirstName = *profileInfo.FirstName
	}
	if profileInfo.LastName != nil {
		user.LastName = *profileInfo.LastName
	}
	if profileInfo.Email != nil && user.Email != *profileInfo.Email {
		userService.enterNewEmail(user.FirstName, user.LastName, *profileInfo.Email, profileInfo.EmailSubject, profileInfo.TemplateFile)
		user.Email = *profileInfo.Email
		user.EmailVerified = false
	}
	if profileInfo.NationalCode != nil {
		user.NationalCode = *profileInfo.NationalCode
	}
	if profileInfo.ProfilePic != nil {
		profilePicPath := userService.constants.S3BucketPath.GetUserProfilePath(profileInfo.UserID, profileInfo.ProfilePic.Filename)
		userService.s3Storage.UploadObject(enum.ProfilePic, profilePicPath, profileInfo.ProfilePic)
		err := userService.s3Storage.DeleteObject(enum.ProfilePic, user.ProfilePicPath)
		if err != nil {
			loggerimpl.GetLogger().Error("unable to delete object", logger.Error("error:", err))
		}
		user.ProfilePicPath = profilePicPath
	}
	err := userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		panic(err)
	}
}

func (userService *UserService) GetAllPermissions() []userdto.PermissionResponse {
	permissions := userService.userRepository.FindAllPermissions(userService.db)
	permissionsResponse := make([]userdto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionsResponse[i] = userdto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Type.Description(),
			Category:    permission.Category.String(),
		}
	}
	return permissionsResponse
}

func (userService *UserService) getRolePermissions(role *entity.Role) []userdto.PermissionResponse {
	if err := userService.userRepository.FindRolePermissions(userService.db, role); err != nil {
		panic(err)
	}
	permissions := make([]userdto.PermissionResponse, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissions[i] = userdto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Type.Description(),
			Category:    permission.Category.String(),
		}
	}
	return permissions
}

func (userService *UserService) GetAllRoles() []userdto.RoleResponse {
	roles := userService.userRepository.FindAllRoles(userService.db)
	rolesResponse := make([]userdto.RoleResponse, len(roles))
	for i, role := range roles {
		rolesResponse[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: userService.getRolePermissions(role),
		}
	}
	return rolesResponse
}

func (userService *UserService) CreateRole(newRoleRequest userdto.NewRoleRequest) {
	_, exist := userService.userRepository.FindRoleByName(userService.db, newRoleRequest.Name)
	if exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Role, userService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
	}
	role := &entity.Role{
		Name: newRoleRequest.Name,
	}
	err := userService.userRepository.CreateRole(userService.db, role)
	if err != nil {
		panic(err)
	}

	existingPermissions := make(map[uint]bool)
	for _, permissionID := range newRoleRequest.PermissionIDs {
		if existingPermissions[permissionID] {
			continue
		}
		permission, exist := userService.userRepository.FindPermissionByID(userService.db, permissionID)
		if !exist {
			notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
			panic(notFoundError)
		}
		if err := userService.userRepository.AssignPermissionToRole(userService.db, role, permission); err != nil {
			panic(err)
		}
		existingPermissions[permissionID] = true
	}
}

func (userService *UserService) GetRoomDetails(roleID uint) userdto.RoleResponse {
	role, exist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		panic(notFoundError)
	}
	return userdto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: userService.getRolePermissions(role),
	}
}

func (userService *UserService) GetRoleOwners(roleID uint) []userdto.CredentialResponse {
	_, exist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		panic(notFoundError)
	}
	users := userService.userRepository.FindUsersByRoleID(userService.db, roleID)
	userCreds := make([]userdto.CredentialResponse, len(users))
	for i, user := range users {
		profilePic := ""
		if user.ProfilePicPath != "" {
			profilePic = userService.s3Storage.GetPresignedURL(enum.ProfilePic, user.ProfilePicPath, 8*time.Hour)
		}
		userCreds[i] = userdto.CredentialResponse{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Phone:      user.Phone,
			Email:      user.Email,
			NationalID: user.NationalCode,
			ProfilePic: profilePic,
			Status:     user.Status.String(),
		}
	}
	return userCreds
}

func (userService *UserService) GetUserRoles(userID uint) []userdto.RoleResponse {
	user, exist := userService.userRepository.FindUserByID(userService.db, userID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		panic(err)
	}
	roles := make([]userdto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: userService.getRolePermissions(&role),
		}
	}
	return roles
}

func (userService *UserService) DeleteRole(roleID uint) {
	_, exist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		panic(notFoundError)
	}
	if err := userService.userRepository.DeleteRole(userService.db, roleID); err != nil {
		panic(err)
	}
}

func (userService *UserService) UpdateRole(newRoleRequest userdto.UpdateRoleRequest) {
	role, exist := userService.userRepository.FindRoleByID(userService.db, newRoleRequest.RoleID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		panic(notFoundError)
	}

	if newRoleRequest.Name != nil {
		role.Name = *newRoleRequest.Name
		if err := userService.userRepository.UpdateRole(userService.db, role); err != nil {
			panic(err)
		}
	}

	existingPermissions := make(map[uint]bool)
	var permissions []entity.Permission
	for _, permissionID := range newRoleRequest.PermissionIDs {
		if existingPermissions[permissionID] {
			continue
		}
		permission, exist := userService.userRepository.FindPermissionByID(userService.db, permissionID)
		if !exist {
			notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
			panic(notFoundError)
		}
		permissions = append(permissions, *permission)
		existingPermissions[permissionID] = true
	}
	if err := userService.userRepository.ReplaceRolePermissions(userService.db, role, permissions); err != nil {
		panic(err)
	}
}

func (userService *UserService) UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) {
	user, exist := userService.userRepository.FindUserByID(userService.db, userRolesRequest.UserID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		panic(notFoundError)
	}

	existingRoles := make(map[uint]bool)
	var roles []entity.Role
	for _, roleID := range userRolesRequest.RoleIDs {
		if existingRoles[roleID] {
			continue
		}
		role, exist := userService.userRepository.FindRoleByID(userService.db, roleID)
		if !exist {
			notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
			panic(notFoundError)
		}
		roles = append(roles, *role)
		existingRoles[roleID] = true
	}
	if err := userService.userRepository.ReplaceUserRoles(userService.db, user, roles); err != nil {
		panic(err)
	}
}
