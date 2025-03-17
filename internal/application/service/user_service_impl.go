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

func (userService *UserService) passwordValidation(password string) []string {
	var errors []string

	userService.validatePasswordTests(&errors, ".{8,}", password, userService.constants.Tag.MinimumLength)
	userService.validatePasswordTests(&errors, "[a-z]", password, userService.constants.Tag.ContainsLowercase)
	userService.validatePasswordTests(&errors, "[A-Z]", password, userService.constants.Tag.ContainsUppercase)
	userService.validatePasswordTests(&errors, "[0-9]", password, userService.constants.Tag.ContainsNumber)
	userService.validatePasswordTests(&errors, "[^\\d\\w]", password, userService.constants.Tag.ContainsSpecialChar)

	return errors
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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

	var passwordErrors exception.ValidationErrors
	passwordErrorTags := userService.passwordValidation(password)

	for _, tag := range passwordErrorTags {
		passwordErrors.Add(userService.constants.Field.Password, tag)
	}
	if len(passwordErrorTags) > 0 {
		return passwordErrors
	}
	return nil
}

func (userService *UserService) Register(registerInfo userdto.BasicRegisterRequest) {
	err := userService.validateRegisterInfo(registerInfo.Phone, registerInfo.Password)
	if err != nil {
		panic(err)
	}
	hashedPassword, err := hashPassword(registerInfo.Password)
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
		Password:      hashedPassword,
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
	var validationErrors exception.ValidationErrors
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, verifyInfo.Phone)
	if !userExist {
		validationErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotRegistered)
		panic(validationErrors)
	}
	if user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err := userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		if exception.IsOTPExpired(err) {
			validationErrors.Add(userService.constants.Field.OTP, userService.constants.Tag.Expired)
		} else if exception.IsInvalidOTP(err) {
			validationErrors.Add(userService.constants.Field.OTP, userService.constants.Tag.Invalid)
		} else {
			panic(err)
		}
		panic(validationErrors)
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
			permissionNames = append(permissionNames, permission.Name)
		}
	}
	return permissionNames
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) userdto.UserInfoResponse {
	user, userExist := userService.userRepository.FindUserByPhone(userService.db, loginInfo.Phone)
	if !userExist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotRegistered)
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
