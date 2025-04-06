package serviceimpl

import (
	"context"
	"regexp"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	cacherepository "github.com/BargheNo/Backend/internal/domain/repository/redis"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants           *bootstrap.Constants
	otpService          service.OTPService
	smsService          service.SMSService
	jwtService          service.JWTService
	userRepository      repository.UserRepository
	userCacheRepository cacherepository.UserCacheRepository
	db                  database.Database
}

func NewUserService(
	constants *bootstrap.Constants,
	otpService service.OTPService,
	smsService service.SMSService,
	jwtService service.JWTService,
	userRepository repository.UserRepository,
	userCacheRepository cacherepository.UserCacheRepository,
	db database.Database,
) *UserService {
	return &UserService{
		constants:           constants,
		otpService:          otpService,
		smsService:          smsService,
		jwtService:          jwtService,
		userRepository:      userRepository,
		userCacheRepository: userCacheRepository,
		db:                  db,
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

func (userService *UserService) validateRegisterInfo(phone, password string) error {
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

	return userService.passwordValidation(password)
}

func (userService *UserService) GetUserCredential(userID uint) userdto.CredentialResponse {
	user, userExist := userService.userRepository.FindUserByID(userService.db, userID)
	if !userExist {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		panic(notFoundError)
	}
	return userdto.CredentialResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}
}

func (userService *UserService) Register(registerInfo userdto.BasicRegisterRequest) {
	err := userService.validateRegisterInfo(registerInfo.Phone, registerInfo.Password)
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

func (userService *UserService) FindUserPermissions(user *entity.User) []string {
	err := userService.userRepository.FindUserRoles(userService.db, user)
	if err != nil {
		panic(err)
	}
	var permissionNames []string
	for _, role := range user.Roles {
		err = userService.userRepository.FindRolePermissions(userService.db, &role)
		if err != nil {
			panic(err)
		}
		permissions := role.Permissions
		for _, permission := range permissions {
			permissionNames = append(permissionNames, permission.Type.String())
		}
	}
	return permissionNames
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
