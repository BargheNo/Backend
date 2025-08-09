package service

import (
	"context"
	"regexp"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/communication"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/repository/redis"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants           *bootstrap.Constants
	otpService          usecase.OTPService
	jwtService          usecase.JWTService
	smsService          communication.SMSService
	emailService        communication.EmailService
	rabbitMQ            message.Broker
	s3Storage           s3.S3Storage
	userRepository      postgres.UserRepository
	userCacheRepository redis.UserCacheRepository
	db                  database.Database
}

type UserServiceDeps struct {
	Constants           *bootstrap.Constants
	OTPService          usecase.OTPService
	JWTService          usecase.JWTService
	SMSService          communication.SMSService
	EmailService        communication.EmailService
	RabbitMQ            message.Broker
	S3Storage           s3.S3Storage
	UserRepository      postgres.UserRepository
	UserCacheRepository redis.UserCacheRepository
	DB                  database.Database
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:           deps.Constants,
		otpService:          deps.OTPService,
		jwtService:          deps.JWTService,
		smsService:          deps.SMSService,
		emailService:        deps.EmailService,
		rabbitMQ:            deps.RabbitMQ,
		s3Storage:           deps.S3Storage,
		userRepository:      deps.UserRepository,
		userCacheRepository: deps.UserCacheRepository,
		db:                  deps.DB,
	}
}

func (userService *UserService) getSortByColumn(requested uint) string {
	allowed := sortby.GetUserSortableColumns()
	sortBy := sortby.UserSortBy(requested)
	if _, ok := allowed[sortBy]; ok {
		return sortBy.DBColumn()
	}
	return sortby.NewsSortByCreatedAt.DBColumn()
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
	data, err := userService.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}
	if data != nil {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, err := userService.userRepository.FindUserByEmail(userService.db, email)
	if err != nil {
		return err
	}
	if user != nil && user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) validateDuplicatePhone(phone string) error {
	var conflictErrors exception.ConflictErrors
	redisKey := userService.constants.RedisKey.GenerateOTPKey(phone)
	data, err := userService.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}
	if data != nil {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return err
	}
	if user != nil && user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) enterNewEmail(firstName, lastName, email, emailSubject, templateFile string) error {
	err := userService.validateDuplicateEmail(email)
	if err != nil {
		return err
	}

	otp, expiryMinute, err := userService.otpService.GenerateOTP()
	if err != nil {
		return err
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(email)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		return err
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
	if err := userService.emailService.SendEmail(email, emailSubject, templateFile, data); err != nil {
		return err
	}
	return nil
}

func (userService *UserService) GetUserSortableColumns() []userdto.UserEnumResponse {
	columns := sortby.GetUserSortableColumns()
	response := make([]userdto.UserEnumResponse, len(columns))
	i := 0
	for col, _ := range columns {
		response[i] = userdto.UserEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
}

func (userService *UserService) IsUserActive(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Status == enum.UserStatusBlock {
		return exception.NewBannedUserForbiddenError()
	}
	return nil
}

func (userService *UserService) GetUserByID(userID uint) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByID(userService.db, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	return user, nil
}

func (userService *UserService) FindActiveUserByPhone(phone string) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return nil, conflictErrors
	}

	return user, nil
}

func (userService *UserService) GetUserCredential(userID uint) (userdto.CredentialResponse, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return userdto.CredentialResponse{}, err
	}

	profilePic := ""
	if user.ProfilePicPath != "" {
		profilePic, err = userService.s3Storage.GetPresignedURL(enum.ProfilePic, user.ProfilePicPath, 8*time.Hour)
		if err != nil {
			return userdto.CredentialResponse{}, err
		}
	}
	return userdto.CredentialResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		Email:      user.Email,
		NationalID: user.NationalCode,
		ProfilePic: profilePic,
		Status:     user.Status.String(),
	}, nil
}

func (userService *UserService) GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	return userService.userRepository.FindUsersByPermission(userService.db, permissionTypes)
}

func (userService *UserService) GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, int64, error) {
	statuses := make([]enum.UserStatus, len(request.Statuses))
	for i, status := range request.Statuses {
		statuses[i] = enum.UserStatus(status)
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(userService.getSortByColumn(request.SortBy), request.Asc)

	users, err := userService.userRepository.FindUserByStatus(userService.db, statuses, options)
	if err != nil {
		return nil, 0, err
	}
	usersResponse := make([]userdto.CredentialResponse, len(users))
	for i, user := range users {
		profilePic := ""
		if user.ProfilePicPath != "" {
			profilePic, err = userService.s3Storage.GetPresignedURL(enum.ProfilePic, user.ProfilePicPath, 8*time.Hour)
			if err != nil {
				return nil, 0, err
			}
		}
		usersResponse[i] = userdto.CredentialResponse{
			ID:         user.ID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Phone:      user.Phone,
			Email:      user.Email,
			NationalID: user.NationalCode,
			ProfilePic: profilePic,
			Status:     user.Status.String(),
		}
	}
	count, err := userService.userRepository.CountUserByStatus(userService.db, statuses)
	if err != nil {
		return nil, 0, err
	}

	return usersResponse, count, nil
}

func (userService *UserService) GetPermissionRoles(request userdto.GetPermissionRolesRequest) ([]userdto.RoleResponse, int64, error) {
	permission, err := userService.userRepository.FindPermissionByID(userService.db, request.PermissionID)
	if err != nil {
		return nil, 0, err
	}
	if permission == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
		return nil, 0, notFoundError
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(userService.getSortByColumn(request.SortBy), request.Asc)

	roles, err := userService.userRepository.FindRolesByPermission(userService.db, request.PermissionID, options)
	if err != nil {
		return nil, 0, err
	}
	rolesResponse := make([]userdto.RoleResponse, len(roles))
	for i, role := range roles {
		rolesResponse[i] = userdto.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}
	count, err := userService.userRepository.CountRolesByPermission(userService.db, request.PermissionID)
	if err != nil {
		return nil, 0, err
	}
	return rolesResponse, count, nil
}

func (userService *UserService) BanUser(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusBlock {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.User, userService.constants.Tag.AlreadyBlocked)
		return conflictErrors
	}
	user.Status = enum.UserStatusBlock
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UnbanUser(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.User, userService.constants.Tag.AlreadyActive)
		return conflictErrors
	}
	user.Status = enum.UserStatusActive
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) Register(registerInfo userdto.BasicRegisterRequest) error {
	err := userService.validateDuplicatePhone(registerInfo.Phone)
	if err != nil {
		return err
	}

	err = userService.passwordValidation(registerInfo.Password)
	if err != nil {
		return err
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), 14)
	if err != nil {
		return err
	}

	err = userService.db.WithTransaction(func(tx database.Database) error {
		err = userService.userRepository.DeleteUserByPhone(tx, registerInfo.Phone)
		if err != nil {
			return err
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
		err = userService.userRepository.CreateUser(tx, user)
		if err != nil {
			return err
		}

		otp, expiryMinute, err := userService.otpService.GenerateOTP()
		if err != nil {
			return err
		}
		redisKey := userService.constants.RedisKey.GenerateOTPKey(registerInfo.Phone)
		err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
		if err != nil {
			return err
		}

		msg := struct {
			UserID uint `json:"userID"`
		}{
			UserID: user.ID,
		}
		if err = userService.rabbitMQ.PublishMessage(userService.constants.RabbitMQ.Events.UserRegistered, msg); err != nil {
			return err
		}
		// userService.smsService.SendOTP(registerInfo.Phone, otp)
		return nil
	})

	return err
}

func (userService *UserService) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error {
	user, err := userService.FindUserByPhone(verifyInfo.Phone)
	if err != nil {
		return err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return notFoundError
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		return err
	}
	user.PhoneVerified = true
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) FindUserPermissions(user *entity.User) ([]userdto.PermissionResponse, error) {
	var permissions []userdto.PermissionResponse
	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		return nil, err
	}
	for _, role := range user.Roles {
		rolePermissions, err := userService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}
	return permissions, nil
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error) {
	user, err := userService.FindActiveUserByPhone(loginInfo.Phone)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		authError := exception.NewInvalidCredentialsError("phone and password not match", nil)
		return userdto.UserInfoResponse{}, authError
	}
	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}
	permissions, err := userService.FindUserPermissions(user)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}
	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Permissions:  permissions,
	}, nil
}

func (userService *UserService) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error {
	_, err := userService.FindActiveUserByPhone(forgotPasswordInfo.Phone)
	if err != nil {
		return err
	}

	otp, expiryMinute, err := userService.otpService.GenerateOTP()
	if err != nil {
		return err
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(forgotPasswordInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		return err
	}
	// userService.smsService.SendOTP(registerInfo.Phone, otp)
	return nil
}

func (userService *UserService) FindUserByPhone(phone string) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	return user, nil
}

func (userService *UserService) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error) {
	user, err := userService.FindUserByPhone(verifyInfo.Phone)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}
	permissions, err := userService.FindUserPermissions(user)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}
	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Permissions:  permissions,
	}, nil
}

func (userService *UserService) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error {
	user, err := userService.GetUserByID(completeRegisterInfo.UserID)
	if err != nil {
		return err
	}

	if completeRegisterInfo.Email != "" {
		userService.enterNewEmail(user.FirstName, user.LastName, completeRegisterInfo.Email, completeRegisterInfo.EmailSubject, completeRegisterInfo.TemplateFile)
	}
	user.Email = completeRegisterInfo.Email
	user.EmailVerified = false
	user.NationalCode = completeRegisterInfo.NationalCode

	err = userService.db.WithTransaction(func(tx database.Database) error {
		if completeRegisterInfo.ProfilePic != nil {
			profilePicPath := userService.constants.S3BucketPath.GetUserProfilePath(completeRegisterInfo.UserID, completeRegisterInfo.ProfilePic.Filename)
			userService.s3Storage.UploadObject(enum.ProfilePic, profilePicPath, completeRegisterInfo.ProfilePic)
			user.ProfilePicPath = profilePicPath
		}
		err = userService.userRepository.UpdateUser(tx, user)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (userService *UserService) VerifyEmail(verifyInfo userdto.VerifyEmailRequest) error {
	var conflictErrors exception.ConflictErrors
	user, err := userService.GetUserByID(verifyInfo.UserID)
	if err != nil {
		return err
	}

	if !user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return conflictErrors
	}
	if user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Email)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		return err
	}
	user.EmailVerified = true
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error {
	user, err := userService.GetUserByID(resetPassInfo.UserID)
	if err != nil {
		return err
	}

	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return conflictErrors
	}

	if err := userService.passwordValidation(resetPassInfo.Password); err != nil {
		return err
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(resetPassInfo.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashesPasswordBytes)

	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UpdateProfile(profileInfo userdto.UpdateProfileRequest) error {
	user, err := userService.GetUserByID(profileInfo.UserID)
	if err != nil {
		return err
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

	oldProfilePicPath := ""
	if profileInfo.ProfilePic != nil {
		profilePicPath := userService.constants.S3BucketPath.GetUserProfilePath(profileInfo.UserID, profileInfo.ProfilePic.Filename)
		userService.s3Storage.UploadObject(enum.ProfilePic, profilePicPath, profileInfo.ProfilePic)
		oldProfilePicPath = user.ProfilePicPath
		user.ProfilePicPath = profilePicPath
	}
	err = userService.db.WithTransaction(func(tx database.Database) error {
		if err := userService.userRepository.UpdateUser(tx, user); err != nil {
			return err
		}

		if oldProfilePicPath != "" {
			if err = userService.s3Storage.DeleteObject(enum.ProfilePic, oldProfilePicPath); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (userService *UserService) GetAllPermissions() ([]userdto.PermissionResponse, error) {
	permissions, err := userService.userRepository.FindAllPermissions(userService.db)
	if err != nil {
		return nil, err
	}
	permissionsResponse := make([]userdto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionsResponse[i] = userdto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Type.Description(),
			Category:    permission.Category.String(),
		}
	}
	return permissionsResponse, nil
}

func (userService *UserService) getRolePermissions(role *entity.Role) ([]userdto.PermissionResponse, error) {
	if err := userService.userRepository.FindRolePermissions(userService.db, role); err != nil {
		return nil, err
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
	return permissions, nil
}

func (userService *UserService) GetAllRoles() ([]userdto.RoleResponse, error) {
	roles, err := userService.userRepository.FindAllRoles(userService.db)
	if err != nil {
		return nil, err
	}
	rolesResponse := make([]userdto.RoleResponse, len(roles))
	for i, role := range roles {
		permissions, err := userService.getRolePermissions(role)
		if err != nil {
			return nil, err
		}
		rolesResponse[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}
	return rolesResponse, nil
}

func (userService *UserService) getPermission(permissionID uint) (*entity.Permission, error) {
	permission, err := userService.userRepository.FindPermissionByID(userService.db, permissionID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Permission}
		return nil, notFoundError
	}
	return permission, nil
}

func (userService *UserService) getRole(roleID uint) (*entity.Role, error) {
	role, err := userService.userRepository.FindRoleByID(userService.db, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.Role}
		return nil, notFoundError
	}
	return role, nil
}

func (userService *UserService) CreateRole(newRoleRequest userdto.NewRoleRequest) error {
	existingRole, err := userService.userRepository.FindRoleByName(userService.db, newRoleRequest.Name)
	if err != nil {
		return err
	}
	if existingRole != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Role, userService.constants.Tag.AlreadyExist)
		return conflictErrors
	}
	err = userService.db.WithTransaction(func(tx database.Database) error {
		role := &entity.Role{
			Name: newRoleRequest.Name,
		}
		err = userService.userRepository.CreateRole(tx, role)
		if err != nil {
			return err
		}

		existingPermissions := make(map[uint]bool)
		for _, permissionID := range newRoleRequest.PermissionIDs {
			if existingPermissions[permissionID] {
				continue
			}

			permission, err := userService.getPermission(permissionID)
			if err != nil {
				return err
			}

			if err := userService.userRepository.AssignPermissionToRole(tx, role, permission); err != nil {
				return err
			}
			existingPermissions[permissionID] = true
		}

		return nil
	})
	return nil
}

func (userService *UserService) GetRoleDetails(roleID uint) (userdto.RoleResponse, error) {
	role, err := userService.getRole(roleID)
	if err != nil {
		return userdto.RoleResponse{}, err
	}

	permissions, err := userService.getRolePermissions(role)
	if err != nil {
		return userdto.RoleResponse{}, err
	}

	return userdto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (userService *UserService) GetRoleOwners(roleID uint) ([]userdto.CredentialResponse, error) {
	_, err := userService.getRole(roleID)
	if err != nil {
		return nil, err
	}

	users, err := userService.userRepository.FindUsersByRoleID(userService.db, roleID)
	if err != nil {
		return nil, err
	}

	userCreds := make([]userdto.CredentialResponse, len(users))
	for i, user := range users {
		profilePic := ""
		if user.ProfilePicPath != "" {
			profilePic, err = userService.s3Storage.GetPresignedURL(enum.ProfilePic, user.ProfilePicPath, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}
		userCreds[i] = userdto.CredentialResponse{
			ID:         user.ID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Phone:      user.Phone,
			Email:      user.Email,
			NationalID: user.NationalCode,
			ProfilePic: profilePic,
			Status:     user.Status.String(),
		}
	}
	return userCreds, nil
}

func (userService *UserService) GetUserRoles(userID uint) ([]userdto.RoleResponse, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if err := userService.userRepository.FindUserRoles(userService.db, user); err != nil {
		return nil, err
	}
	roles := make([]userdto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		permissions, err := userService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		roles[i] = userdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}
	return roles, nil
}

func (userService *UserService) DeleteRole(roleID uint) error {
	_, err := userService.getRole(roleID)
	if err != nil {
		return err
	}

	if err := userService.userRepository.DeleteRole(userService.db, roleID); err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UpdateRole(newRoleRequest userdto.UpdateRoleRequest) error {
	role, err := userService.getRole(newRoleRequest.RoleID)
	if err != nil {
		return err
	}

	existingPermissions := make(map[uint]bool)
	var permissions []entity.Permission
	for _, permissionID := range newRoleRequest.PermissionIDs {
		if existingPermissions[permissionID] {
			continue
		}

		permission, err := userService.getPermission(permissionID)
		if err != nil {
			return err
		}

		permissions = append(permissions, *permission)
		existingPermissions[permissionID] = true
	}

	err = userService.db.WithTransaction(func(tx database.Database) error {
		if newRoleRequest.Name != nil {
			role.Name = *newRoleRequest.Name
			if err := userService.userRepository.UpdateRole(tx, role); err != nil {
				return err
			}
		}

		if err := userService.userRepository.ReplaceRolePermissions(tx, role, permissions); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (userService *UserService) UpdateUserRoles(userRolesRequest userdto.UpdateUserRolesRequest) error {
	user, err := userService.GetUserByID(userRolesRequest.UserID)
	if err != nil {
		return err
	}

	existingRoles := make(map[uint]bool)
	var roles []entity.Role
	for _, roleID := range userRolesRequest.RoleIDs {
		if existingRoles[roleID] {
			continue
		}

		role, err := userService.getRole(roleID)
		if err != nil {
			return err
		}

		roles = append(roles, *role)
		existingRoles[roleID] = true
	}

	if err := userService.userRepository.ReplaceUserRoles(userService.db, user, roles); err != nil {
		return err
	}

	return nil
}
