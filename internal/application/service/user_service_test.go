package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"

	databaseMocks "github.com/BargheNo/Backend/mocks/infrastructure/database"
	brokerMocks "github.com/BargheNo/Backend/mocks/infrastructure/rabbitmq"
	recaptchaMocks "github.com/BargheNo/Backend/mocks/infrastructure/recaptcha"
	repositoryMocks "github.com/BargheNo/Backend/mocks/infrastructure/repository/postgres"
	cacheMocks "github.com/BargheNo/Backend/mocks/infrastructure/repository/redis"
	s3Mocks "github.com/BargheNo/Backend/mocks/infrastructure/s3"
	serviceMocks "github.com/BargheNo/Backend/mocks/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	constants           *bootstrap.Constants
	otpService          *serviceMocks.OTPServiceMock
	jwtService          *serviceMocks.JWTServiceMock
	smsService          *serviceMocks.SMSServiceMock
	emailService        *serviceMocks.EmailServiceMock
	rbacService         *serviceMocks.RBACServiceMock
	rabbitMQ            *brokerMocks.BrokerMock
	s3Storage           *s3Mocks.S3StorageMock
	userRepository      *repositoryMocks.UserRepositoryMock
	userCacheRepository *cacheMocks.UserCacheRepositoryMock
	db                  *databaseMocks.DatabaseMock
	recaptcha           *recaptchaMocks.RecaptchaMock
	userService         *UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	suite.constants = config.Constants
	suite.otpService = serviceMocks.NewOTPServiceMock()
	suite.jwtService = serviceMocks.NewJWTServiceMock()
	suite.smsService = serviceMocks.NewSMSServiceMock()
	suite.emailService = serviceMocks.NewEmailServiceMock()
	suite.rbacService = serviceMocks.NewRBACServiceMock()
	suite.rabbitMQ = brokerMocks.NewBrokerMock()
	suite.s3Storage = s3Mocks.NewS3StorageMock()
	suite.userRepository = repositoryMocks.NewUserRepositoryMock()
	suite.userCacheRepository = cacheMocks.NewUserCacheRepositoryMock()
	suite.db = databaseMocks.NewDatabaseMock()
	suite.recaptcha = recaptchaMocks.NewRecaptchaMock()

	deps := UserServiceDeps{
		Constants:           suite.constants,
		OTPService:          suite.otpService,
		JWTService:          suite.jwtService,
		SMSService:          suite.smsService,
		EmailService:        suite.emailService,
		RBACService:         suite.rbacService,
		RabbitMQ:            suite.rabbitMQ,
		S3Storage:           suite.s3Storage,
		UserRepository:      suite.userRepository,
		UserCacheRepository: suite.userCacheRepository,
		DB:                  suite.db,
		Recaptcha:           suite.recaptcha,
	}

	suite.userService = NewUserService(deps)
}

func (s *UserServiceTestSuite) TestEnterNewEmail() {
	s.Run("success - Email is not registered", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, nil).Once()

		s.otpService.On("GenerateOTP").Return("123456", 1, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Duplicate Email", func() {
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, nil).Once()

		err := s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")

		s.Error(err)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
	})
	s.Run("Error - Generate OTP Error", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, nil).Once()
		s.otpService.On("GenerateOTP").Return("", 0, errors.New("generate OTP error")).Once()

		err := s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")

		s.Error(err)

		s.userCacheRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("Error - Set OTP to Cache Error", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, nil).Once()
		s.otpService.On("GenerateOTP").Return("123456", 1, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(errors.New("set OTP to cache error")).Once()

		err := s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")

		s.Error(err)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Send Email Error", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, nil).Once()

		s.otpService.On("GenerateOTP").Return("123456", 10, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("send email error")).Once()

		err := s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")

		s.Error(err)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserSortableColumns() {
	s.Run("success - returns user sortable columns", func() {
		result := s.userService.GetUserSortableColumns()

		s.NotNil(result)
		s.Greater(len(result), 0)

		// Verify each column has ID and Name
		for _, col := range result {
			s.Greater(col.ID, uint(0))
			s.NotEmpty(col.Name)
		}
	})
}

func (s *UserServiceTestSuite) TestIsUserActive() {
	s.Run("success - user is active", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusActive,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		err := s.userService.IsUserActive(userID)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user is blocked", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusBlock,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		err := s.userService.IsUserActive(userID)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		var nilUser *entity.User = nil
		userID := uint(1)

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, nil).Once()

		err := s.userService.IsUserActive(userID)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserByID() {
	s.Run("success - user found", func() {
		userID := uint(1)
		user := &entity.User{
			FirstName: "John",
			LastName:  "Doe",
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		result, err := s.userService.GetUserByID(userID)

		s.NoError(err)
		s.Equal(user, result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		var nilUser *entity.User = nil
		userID := uint(1)

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, nil).Once()

		result, err := s.userService.GetUserByID(userID)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		var nilUser *entity.User = nil
		userID := uint(1)

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, errors.New("database error")).Once()

		result, err := s.userService.GetUserByID(userID)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestFindActiveUserByPhone() {
	s.Run("success - active user found", func() {
		phone := "+1234567890"
		user := &entity.User{
			Phone:         phone,
			PhoneVerified: true,
		}

		s.userRepository.On("FindUserByPhone", s.db, phone).Return(user, nil).Once()

		result, err := s.userService.FindActiveUserByPhone(phone)

		s.NoError(err)
		s.Equal(user, result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		phone := "+1234567890"
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByPhone", s.db, phone).Return(nilUser, nil).Once()

		result, err := s.userService.FindActiveUserByPhone(phone)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - phone not verified", func() {
		phone := "+1234567890"
		user := &entity.User{
			Phone:         phone,
			PhoneVerified: false,
		}

		s.userRepository.On("FindUserByPhone", s.db, phone).Return(user, nil).Once()

		result, err := s.userService.FindActiveUserByPhone(phone)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserCredential() {
	s.Run("success - user credential with profile pic", func() {
		userID := uint(1)
		user := &entity.User{
			FirstName:      "John",
			LastName:       "Doe",
			Phone:          "+1234567890",
			Email:          "john@example.com",
			EmailVerified:  true,
			NationalCode:   "1234567890",
			ProfilePicPath: "profile/pic.jpg",
			Status:         enum.UserStatusActive,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.s3Storage.On("GetPresignedURL", enum.ProfilePic, "profile/pic.jpg", 8*time.Hour).Return("https://presigned-url.com", nil).Once()

		_, err := s.userService.GetUserCredential(userID)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})

	s.Run("success - user credential without profile pic", func() {
		userID := uint(1)
		user := &entity.User{
			FirstName:      "John",
			LastName:       "Doe",
			Phone:          "+1234567890",
			Email:          "john@example.com",
			EmailVerified:  true,
			NationalCode:   "1234567890",
			ProfilePicPath: "",
			Status:         enum.UserStatusActive,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		_, err := s.userService.GetUserCredential(userID)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - s3 storage error", func() {
		userID := uint(1)
		user := &entity.User{
			FirstName:      "John",
			LastName:       "Doe",
			ProfilePicPath: "profile/pic.jpg",
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.s3Storage.On("GetPresignedURL", enum.ProfilePic, "profile/pic.jpg", 8*time.Hour).Return("", errors.New("s3 error")).Once()

		result, err := s.userService.GetUserCredential(userID)

		s.Error(err)
		s.Equal(userdto.CredentialResponse{}, result)
		s.userRepository.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUsersByPermission() {
	s.Run("success - users found", func() {
		permissions := []enum.PermissionType{enum.PermissionAll}
		users := []*entity.User{
			{FirstName: "John"},
			{FirstName: "Jane"},
		}

		s.userRepository.On("FindUsersByPermission", s.db, permissions).Return(users, nil).Once()

		result, err := s.userService.GetUsersByPermission(permissions)

		s.NoError(err)
		s.Equal(users, result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		permissions := []enum.PermissionType{enum.PermissionAll}
		var nilUsers []*entity.User = nil

		s.userRepository.On("FindUsersByPermission", s.db, permissions).Return(nilUsers, errors.New("database error")).Once()

		result, err := s.userService.GetUsersByPermission(permissions)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUsersByStatus() {
	s.Run("success - users found with query", func() {
		request := userdto.GetUsersListRequest{
			Query:  "John",
			Status: uint(enum.UserStatusActive),
			Limit:  10,
			Offset: 0,
			SortBy: 1,
			Asc:    true,
		}

		users := []*entity.User{
			{FirstName: "John", LastName: "Doe", Phone: "+1234567890", Email: "john@example.com", EmailVerified: true, NationalCode: "1234567890", Status: enum.UserStatusActive},
		}

		s.userRepository.On("FindUsersByQuery", s.db, "John", mock.Anything).Return(users, nil).Once()
		s.userRepository.On("CountUsersByQuery", s.db, "John").Return(int64(1), nil).Once()

		result, count, err := s.userService.GetUsersByStatus(request)

		s.NoError(err)
		s.Equal(int64(1), count)
		s.Len(result, 1)
		s.Equal("John", result[0].FirstName)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("success - users found without query", func() {
		request := userdto.GetUsersListRequest{
			Status: uint(enum.UserStatusActive),
			Limit:  10,
			Offset: 0,
			SortBy: 1,
			Asc:    true,
		}

		users := []*entity.User{
			{FirstName: "John", LastName: "Doe", Phone: "+1234567890", Email: "john@example.com", EmailVerified: true, NationalCode: "1234567890", Status: enum.UserStatusActive},
		}

		s.userRepository.On("FindUserByStatus", s.db, []enum.UserStatus{enum.UserStatusActive}, mock.Anything).Return(users, nil).Once()
		s.userRepository.On("CountUserByStatus", s.db, []enum.UserStatus{enum.UserStatusActive}).Return(int64(1), nil).Once()

		result, count, err := s.userService.GetUsersByStatus(request)

		s.NoError(err)
		s.Equal(int64(1), count)
		s.Len(result, 1)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestBanUser() {
	s.Run("success - user banned", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusActive,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		err := s.userService.BanUser(userID)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user already blocked", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusBlock,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		err := s.userService.BanUser(userID)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUnbanUser() {
	s.Run("success - user unbanned", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusBlock,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		err := s.userService.UnbanUser(userID)

		s.NoError(err)
		s.Equal(enum.UserStatusActive, user.Status)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user already active", func() {
		userID := uint(1)
		user := &entity.User{
			Status: enum.UserStatusActive,
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()

		err := s.userService.UnbanUser(userID)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestRegister() {
	s.Run("success - user registered", func() {
		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+1234567890",
			Password:  "Password123!",
			Recaptcha: "recaptcha-token",
		}

		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.recaptcha.On("Verify", "recaptcha-token").Return(nil).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(nilUser, nil).Once()
		s.userRepository.On("DeleteUserByPhone", mock.Anything, "+1234567890").Return(nil).Once()
		s.userRepository.On("CreateUser", mock.Anything, mock.Anything).Return(nil).Once()
		s.otpService.On("GenerateOTP").Return("123456", 10, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.rabbitMQ.On("PublishMessage", mock.Anything, mock.Anything).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.userService.Register(request)

		s.NoError(err)
		s.recaptcha.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.rabbitMQ.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - recaptcha verification failed", func() {
		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+1234567890",
			Password:  "Password123!",
			Recaptcha: "invalid-token",
		}

		s.recaptcha.On("Verify", "invalid-token").Return(errors.New("recaptcha error")).Once()

		err := s.userService.Register(request)

		s.Error(err)
		s.recaptcha.AssertExpectations(s.T())
	})

	s.Run("error - password validation failed", func() {
		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+1234567890",
			Password:  "weak",
			Recaptcha: "recaptcha-token",
		}

		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.recaptcha.On("Verify", "recaptcha-token").Return(nil).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(nilUser, nil).Once()

		err := s.userService.Register(request)

		s.Error(err)
		s.recaptcha.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyPhone() {
	s.Run("success - phone verified", func() {
		request := userdto.VerifyPhoneRequest{
			Phone: "+1234567890",
			OTP:   "123456",
		}

		user := &entity.User{
			Phone:         "+1234567890",
			PhoneVerified: false,
		}

		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(user, nil).Once()
		s.otpService.On("VerifyOTP", mock.Anything, "123456").Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		err := s.userService.VerifyPhone(request)

		s.NoError(err)
		s.True(user.PhoneVerified)
		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		request := userdto.VerifyPhoneRequest{
			Phone: "+1234567890",
			OTP:   "123456",
		}

		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(nilUser, nil).Once()

		err := s.userService.VerifyPhone(request)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestLogin() {
	s.Run("success - login successful", func() {
		request := userdto.LoginRequest{
			Phone:     "+1234567890",
			Password:  "Password123!",
			Recaptcha: "recaptcha-token",
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
		user := &entity.User{
			Phone:         "+1234567890",
			PhoneVerified: true,
			Password:      string(hashed),
			FirstName:     "John",
			LastName:      "Doe",
		}

		permissions := []rbacdto.PermissionResponse{
			{
				ID:   1,
				Name: "read",
			},
		}

		s.recaptcha.On("Verify", "recaptcha-token").Return(nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(user, nil).Once()
		s.jwtService.On("GenerateToken", mock.Anything).Return("access-token", "refresh-token", nil).Once()
		s.rbacService.On("GetUserPermissions", user).Return(permissions, nil).Once()

		result, err := s.userService.Login(request)

		s.NoError(err)
		s.Equal("access-token", result.AccessToken)
		s.Equal("refresh-token", result.RefreshToken)
		s.Equal("John", result.FirstName)
		s.Equal("Doe", result.LastName)
		s.Equal(permissions, result.Permissions)
		s.recaptcha.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})

	s.Run("error - invalid credentials", func() {
		request := userdto.LoginRequest{
			Phone:     "+1234567890",
			Password:  "WrongPassword123!",
			Recaptcha: "recaptcha-token",
		}

		user := &entity.User{
			Phone:         "+1234567890",
			PhoneVerified: true,
			Password:      "$2a$14$hashedpassword",
		}

		s.recaptcha.On("Verify", "recaptcha-token").Return(nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(user, nil).Once()

		result, err := s.userService.Login(request)

		s.Error(err)
		s.Equal(userdto.UserInfoResponse{}, result)
		s.recaptcha.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestForgotPassword() {
	s.Run("success - OTP sent", func() {
		request := userdto.ForgotPasswordRequest{
			Phone: "+1234567890",
		}

		user := &entity.User{
			Phone:         "+1234567890",
			PhoneVerified: true,
		}

		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(user, nil).Once()
		s.otpService.On("GenerateOTP").Return("123456", 10, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()

		err := s.userService.ForgotPassword(request)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		request := userdto.ForgotPasswordRequest{
			Phone: "+1234567890",
		}

		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(nilUser, nil).Once()

		err := s.userService.ForgotPassword(request)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestFindUserByPhone() {
	s.Run("success - user found", func() {
		phone := "+1234567890"
		user := &entity.User{
			Phone: phone,
		}

		s.userRepository.On("FindUserByPhone", s.db, phone).Return(user, nil).Once()

		result, err := s.userService.FindUserByPhone(phone)

		s.NoError(err)
		s.Equal(user, result)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		phone := "+1234567890"
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByPhone", s.db, phone).Return(nilUser, nil).Once()

		result, err := s.userService.FindUserByPhone(phone)

		s.Error(err)
		s.Nil(result)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyOTP() {
	s.Run("success - OTP verified", func() {
		request := userdto.VerifyPhoneRequest{
			Phone: "+1234567890",
			OTP:   "123456",
		}

		user := &entity.User{
			Phone:         "+1234567890",
			PhoneVerified: true,
		}

		permissions := []rbacdto.PermissionResponse{
			{
				ID:   1,
				Name: "read",
			},
		}

		s.userRepository.On("FindUserByPhone", s.db, "+1234567890").Return(user, nil).Once()
		s.otpService.On("VerifyOTP", mock.Anything, "123456").Return(nil).Once()
		s.jwtService.On("GenerateToken", mock.Anything).Return("access-token", "refresh-token", nil).Once()
		s.rbacService.On("GetUserPermissions", user).Return(permissions, nil).Once()

		result, err := s.userService.VerifyOTP(request)

		s.NoError(err)
		s.Equal("access-token", result.AccessToken)
		s.Equal("refresh-token", result.RefreshToken)
		s.Equal(permissions, result.Permissions)
		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestCompleteRegister() {
	s.Run("success - registration completed", func() {
		request := userdto.CompleteRegisterRequest{
			UserID:       1,
			Email:        "john@example.com",
			EmailSubject: "Verify Email",
			TemplateFile: "template.html",
			NationalCode: "1234567890",
		}

		user := &entity.User{
			FirstName: "John",
			LastName:  "Doe",
		}

		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, "john@example.com").Return(nilUser, nil).Once()

		s.otpService.On("GenerateOTP").Return("123456", 5, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		s.userRepository.On("UpdateUser", mock.Anything, mock.Anything).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.userService.CompleteRegister(request)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyEmail() {
	s.Run("success - email verified", func() {
		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "john@example.com",
			OTP:    "123456",
		}

		user := &entity.User{
			PhoneVerified: true,
			EmailVerified: false,
		}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
		s.otpService.On("VerifyOTP", mock.Anything, "123456").Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		err := s.userService.VerifyEmail(request)

		s.NoError(err)
		s.True(user.EmailVerified)
		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})

	s.Run("error - phone not verified", func() {
		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "john@example.com",
			OTP:    "123456",
		}

		user := &entity.User{
			PhoneVerified: false,
			EmailVerified: false,
		}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()

		err := s.userService.VerifyEmail(request)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestResetPassword() {
	s.Run("success - password reset", func() {
		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword123!",
		}

		user := &entity.User{
			PhoneVerified: true,
		}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		err := s.userService.ResetPassword(request)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - phone not verified", func() {
		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword123!",
		}

		user := &entity.User{
			PhoneVerified: false,
		}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()

		err := s.userService.ResetPassword(request)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - invalid password", func() {
		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "weak",
		}

		user := &entity.User{
			PhoneVerified: true,
		}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()

		err := s.userService.ResetPassword(request)

		s.Error(err)
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUpdateProfile() {
	s.Run("success - profile updated", func() {
		request := userdto.UpdateProfileRequest{
			UserID:       1,
			FirstName:    stringPtr("Jane"),
			LastName:     stringPtr("Smith"),
			Email:        stringPtr("jane@example.com"),
			NationalCode: stringPtr("0987654321"),
		}

		user := &entity.User{
			FirstName: "John",
			LastName:  "Doe",
		}

		var nilUser *entity.User

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nil, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, "jane@example.com").Return(nilUser, nil).Once()
		s.otpService.On("GenerateOTP").Return("123456", 5, nil).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.userService.UpdateProfile(request)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserRoles() {
	s.Run("success - user roles retrieved", func() {
		userID := uint(1)
		user := &entity.User{}

		roles := []rbacdto.RoleResponse{
			{ID: 1, Name: "Admin"},
			{ID: 2, Name: "User"},
		}

		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.rbacService.On("GetUserRoles", user).Return(roles, nil).Once()

		result, err := s.userService.GetUserRoles(userID)

		s.NoError(err)
		s.Equal(roles, result)
		s.userRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUpdateUserRoles() {
	s.Run("success - user roles updated", func() {
		request := userdto.UpdateUserRolesRequest{
			UserID:  1,
			RoleIDs: []uint{1, 2},
		}

		user := &entity.User{}

		s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
		s.rbacService.On("UpdateUserRoles", mock.Anything).Return(nil).Once()

		err := s.userService.UpdateUserRoles(request)

		s.NoError(err)
		s.userRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetRoleOwners() {
	s.Run("success - role owners retrieved", func() {
		request := rbacdto.GetRoleOwnersRequest{
			RoleID: 1,
		}

		user := &entity.User{}

		s.rbacService.On("GetRoleOwners", request).Return([]*entity.User{user}, int64(1), nil).Once()
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, nil).Once()

		_, _, err := s.userService.GetRoleOwners(request)

		s.NoError(err)
		s.rbacService.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestRefreshToken() {
	s.Run("success - token refreshed", func() {
		refreshToken := "refresh-token"
		userID := uint(1)

		claims := map[string]interface{}{
			"sub": float64(userID),
		}

		user := &entity.User{
			FirstName: "John",
			LastName:  "Doe",
		}

		permissions := []rbacdto.PermissionResponse{
			{
				ID:   1,
				Name: "read",
			},
		}

		s.jwtService.On("ValidateToken", refreshToken).Return(claims, nil).Once()
		s.userRepository.On("FindUserByID", s.db, userID).Return(user, nil).Once()
		s.jwtService.On("GenerateToken", userID).Return("new-access-token", "new-refresh-token", nil).Once()
		s.rbacService.On("GetUserPermissions", user).Return(permissions, nil).Once()

		result, err := s.userService.RefreshToken(refreshToken)

		s.NoError(err)
		s.Equal("new-access-token", result.AccessToken)
		s.Equal(refreshToken, result.RefreshToken)
		s.Equal("John", result.FirstName)
		s.Equal("Doe", result.LastName)
		s.Equal(permissions, result.Permissions)
		s.jwtService.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})

	s.Run("error - invalid refresh token", func() {
		refreshToken := "invalid-token"

		s.jwtService.On("ValidateToken", refreshToken).Return(nil, errors.New("invalid token")).Once()

		result, err := s.userService.RefreshToken(refreshToken)

		s.Error(err)
		s.Equal(userdto.UserInfoResponse{}, result)
		s.jwtService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestPasswordValidation() {
	s.Run("success - valid password", func() {
		password := "ValidPassword123!"

		err := s.userService.passwordValidation(password)

		s.NoError(err)
	})

	s.Run("error - password too short", func() {
		password := "Short1!"

		err := s.userService.passwordValidation(password)

		s.Error(err)
	})

	s.Run("error - missing lowercase", func() {
		password := "UPPERCASE123!"

		err := s.userService.passwordValidation(password)

		s.Error(err)
	})

	s.Run("error - missing uppercase", func() {
		password := "lowercase123!"

		err := s.userService.passwordValidation(password)

		s.Error(err)
	})

	s.Run("error - missing number", func() {
		password := "NoNumbers!"

		err := s.userService.passwordValidation(password)

		s.Error(err)
	})

	s.Run("error - missing special character", func() {
		password := "NoSpecialChar123"

		err := s.userService.passwordValidation(password)

		s.Error(err)
	})
}

func (s *UserServiceTestSuite) TestValidateDuplicateEmail() {
	s.Run("success - email not duplicate", func() {
		email := "new@example.com"
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, email).Return(nilUser, nil).Once()

		err := s.userService.validateDuplicateEmail(email)

		s.NoError(err)
		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - email already in cache", func() {
		email := "existing@example.com"
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, nil).Once()

		err := s.userService.validateDuplicateEmail(email)

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
	})

	s.Run("error - email already registered", func() {
		email := "existing@example.com"
		user := &entity.User{
			Email:         email,
			EmailVerified: true,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nil, nil).Once()
		s.userRepository.On("FindUserByEmail", s.db, email).Return(user, nil).Once()

		err := s.userService.validateDuplicateEmail(email)

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestValidateDuplicatePhone() {
	s.Run("success - phone not duplicate", func() {
		phone := "+1234567890"
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, phone).Return(nilUser, nil).Once()

		err := s.userService.validateDuplicatePhone(phone)

		s.NoError(err)
		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("error - phone already in cache", func() {
		phone := "+1234567890"
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, nil).Once()

		err := s.userService.validateDuplicatePhone(phone)

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
	})

	s.Run("error - phone already registered", func() {
		phone := "+1234567890"
		user := &entity.User{
			Phone:         phone,
			PhoneVerified: true,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nil, nil).Once()
		s.userRepository.On("FindUserByPhone", s.db, phone).Return(user, nil).Once()

		err := s.userService.validateDuplicatePhone(phone)

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func stringPtr(s string) *string {
	return &s
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
