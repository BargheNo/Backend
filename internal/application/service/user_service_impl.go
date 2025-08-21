package service

import (
	"context"
	"regexp"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/communication"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/recaptcha"
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
	rbacService         usecase.RBACService
	rabbitMQ            message.Broker
	s3Storage           s3.S3Storage
	userRepository      postgres.UserRepository
	userCacheRepository redis.UserCacheRepository
	db                  database.Database
	recaptcha           recaptcha.Recaptcha
}

type UserServiceDeps struct {
	Constants           *bootstrap.Constants
	OTPService          usecase.OTPService
	JWTService          usecase.JWTService
	SMSService          communication.SMSService
	EmailService        communication.EmailService
	RBACService         usecase.RBACService
	RabbitMQ            message.Broker
	S3Storage           s3.S3Storage
	UserRepository      postgres.UserRepository
	UserCacheRepository redis.UserCacheRepository
	DB                  database.Database
	Recaptcha           recaptcha.Recaptcha
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:           deps.Constants,
		otpService:          deps.OTPService,
		jwtService:          deps.JWTService,
		smsService:          deps.SMSService,
		emailService:        deps.EmailService,
		rbacService:         deps.RBACService,
		rabbitMQ:            deps.RabbitMQ,
		s3Storage:           deps.S3Storage,
		userRepository:      deps.UserRepository,
		userCacheRepository: deps.UserCacheRepository,
		db:                  deps.DB,
		recaptcha:           deps.Recaptcha,
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
	if user == nil || !user.PhoneVerified {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
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
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		NationalID:    user.NationalCode,
		ProfilePic:    profilePic,
		Status:        user.Status.String(),
	}, nil
}

func (userService *UserService) GetUsersByPermission(permissionTypes []enum.PermissionType) ([]*entity.User, error) {
	return userService.userRepository.FindUsersByPermission(userService.db, permissionTypes)
}

func (userService *UserService) mapToFilterStatuses(enumStatus uint) []enum.UserStatus {
	statuses := enum.GetAllUserStatus()
	for _, status := range statuses {
		if uint(status) == enumStatus {
			if status == enum.UserStatusAll {
				return statuses
			}
			return []enum.UserStatus{status}
		}
	}
	return statuses
}

func (userService *UserService) getUsersByQuery(query string, statuses []enum.UserStatus, options *postgres.QueryOptions) ([]*entity.User, int64, error) {
	if query == "" {
		users, err := userService.userRepository.FindUserByStatus(userService.db, statuses, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := userService.userRepository.CountUserByStatus(userService.db, statuses)
		if err != nil {
			return nil, 0, err
		}
		return users, count, nil
	}
	users, err := userService.userRepository.FindUsersByQuery(userService.db, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := userService.userRepository.CountUsersByQuery(userService.db, query)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (userService *UserService) GetUsersByStatus(request userdto.GetUsersListRequest) ([]userdto.CredentialResponse, int64, error) {
	statuses := userService.mapToFilterStatuses(request.Status)

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(userService.getSortByColumn(request.SortBy), request.Asc)

	users, count, err := userService.getUsersByQuery(request.Query, statuses, options)
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
			ID:            user.ID,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Phone:         user.Phone,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			NationalID:    user.NationalCode,
			ProfilePic:    profilePic,
			Status:        user.Status.String(),
		}
	}

	return usersResponse, count, nil
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
	err := userService.recaptcha.Verify(registerInfo.Recaptcha)
	if err != nil {
		return err
	}

	err = userService.validateDuplicatePhone(registerInfo.Phone)
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

func (userService *UserService) RefreshToken(refreshToken string) (userdto.UserInfoResponse, error) {
	claims, err := userService.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	userID := uint(claims["sub"].(float64))
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	accessToken, _, err := userService.jwtService.GenerateToken(userID)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	permissions, err := userService.rbacService.GetUserPermissions(user)
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

func (userService *UserService) Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error) {
	err := userService.recaptcha.Verify(loginInfo.Recaptcha)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

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
	permissions, err := userService.rbacService.GetUserPermissions(user)
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
	permissions, err := userService.rbacService.GetUserPermissions(user)
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

func (userService *UserService) GetUserRoles(userID uint) ([]rbacdto.RoleResponse, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	roles, err := userService.rbacService.GetUserRoles(user)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (userService *UserService) UpdateUserRoles(request userdto.UpdateUserRolesRequest) error {
	user, err := userService.GetUserByID(request.UserID)
	if err != nil {
		return err
	}

	updateRoleRequest := rbacdto.UpdateUserRolesRequest{
		User:    user,
		RoleIDs: request.RoleIDs,
	}
	if err := userService.rbacService.UpdateUserRoles(updateRoleRequest); err != nil {
		return err
	}

	return nil
}

func (userService *UserService) GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]userdto.CredentialResponse, int64, error) {
	users, count, err := userService.rbacService.GetRoleOwners(request)
	if err != nil {
		return nil, 0, err
	}
	userCreds := make([]userdto.CredentialResponse, len(users))
	for i, user := range users {
		cred, err := userService.GetUserCredential(user.ID)
		if err != nil {
			return nil, 0, err
		}
		userCreds[i] = cred
	}

	return userCreds, count, nil
}
