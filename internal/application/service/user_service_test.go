package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	constants           *bootstrap.Constants
	otpService          *mocks.OtpServiceMock
	jwtService          *mocks.JwtServiceMock
	smsService          *mocks.SMSServiceMock
	emailService        *mocks.EmailServiceMock
	s3Storage           *mocks.S3StorageMock
	userRepository      *mocks.UserRepositoryMock
	userCacheRepository *mocks.UserCacheRepositoryMock
	db                  *mocks.DatabaseMock
	userService         *UserService
	rabbitMQ            *mocks.BrokerMock
}

func (s *UserServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	s.constants = config.Constants
	s.otpService = mocks.NewOtpServiceMock()
	s.jwtService = mocks.NewJwtServiceMock()
	s.smsService = mocks.NewSMSServiceMock()
	s.emailService = mocks.NewEmailServiceMock()
	s.s3Storage = mocks.NewS3StorageMock()
	s.userRepository = mocks.NewUserRepositoryMock()
	s.userCacheRepository = mocks.NewUserCacheRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.rabbitMQ = mocks.NewBrokerMock()

	deps := UserServiceDeps{
		Constants:           s.constants,
		OTPService:          s.otpService,
		JWTService:          s.jwtService,
		SMSService:          s.smsService,
		EmailService:        s.emailService,
		S3Storage:           s.s3Storage,
		UserRepository:      s.userRepository,
		UserCacheRepository: s.userCacheRepository,
		DB:                  s.db,
		RabbitMQ:            s.rabbitMQ,
	}
	s.userService = NewUserService(deps)
}

func (s *UserServiceTestSuite) TestValidatePasswordTests() {
	s.Run("success - Password is valid", func() {
		s.userService.validatePasswordTests(nil, "Password@123", "Password@123", "Password@123")
	})
}

func (s *UserServiceTestSuite) TestPasswordValidation() {
	s.Run("success - Password is valid", func() {
		s.userService.passwordValidation("Password@123")
	})
}

func (s *UserServiceTestSuite) TestValidateDuplicateEmail() {
	s.Run("success - Email is not registered", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()

		s.userService.validateDuplicateEmail("test@example.com")

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - Email is registered but not verified", func() {
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()

		response := s.userService.validateDuplicateEmail("test@example.com")
		s.IsType(exception.ConflictErrors{}, response)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - Email is registered and verified", func() {
		user := &entity.User{
			EmailVerified: true,
		}
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(user, true).Once()

		response := s.userService.validateDuplicateEmail("test@example.com")
		s.IsType(exception.ConflictErrors{}, response)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})

}

func (s *UserServiceTestSuite) TestValidateDuplicatePhone() {
	s.Run("success - Phone is not registered", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		s.userService.validateDuplicatePhone("1234567890")

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - Phone is registered but not verified", func() {
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()

		response := s.userService.validateDuplicatePhone("1234567890")
		s.IsType(exception.ConflictErrors{}, response)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - Phone is registered and verified", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		response := s.userService.validateDuplicatePhone("1234567890")
		s.IsType(exception.ConflictErrors{}, response)

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestEnterNewEmail() {
	s.Run("success - Email is not registered", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()

		s.otpService.On("GenerateOTP").Return("123456", 10).Once()
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
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()

		s.Panics(func() {
			s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")
		})

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
	})
	s.Run("Error - Set OTP to Cache Error", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()
		s.otpService.On("GenerateOTP").Return("123456", 10).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(errors.New("set OTP to cache error")).Once()

		s.Panics(func() {
			s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")
		})

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Send Email Error", func() {
		var nilUser *entity.User = nil
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()

		s.otpService.On("GenerateOTP").Return("123456", 10).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, "123456", mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("send email error")).Once()

		s.Panics(func() {
			s.userService.enterNewEmail("John", "Doe", "test@example.com", "test subject", "test template")
		})

		s.userCacheRepository.AssertExpectations(s.T())
		s.userRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
	})

}

func (s *UserServiceTestSuite) TestIsUserActive() {
	s.Run("success - User is active", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()

		s.userService.IsUserActive(userID)

		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("Error - User is not active", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.IsUserActive(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserByID() {
	s.Run("success - User found", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()

		s.userService.GetUserByID(userID)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - User not found", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.GetUserByID(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserCredential() {
	s.Run("success - User Credentials found", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()

		s.userService.GetUserCredential(userID)

		s.userRepository.AssertExpectations(s.T())

	})
	s.Run("Error - User Not Found", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.GetUserCredential(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - Get User With Profile Picture", func() {
		userID := uint(1)
		profilePicPath := "profile.jpg"
		profilePic := "https://example.com/profile.jpg"
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{
			ProfilePicPath: profilePicPath,
		}, true).Once()

		s.s3Storage.On("GetPresignedURL", enum.ProfilePic, profilePicPath, 8*time.Hour).Return(profilePic).Once()

		s.userService.GetUserCredential(userID)

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUsersByPermission() {
	s.Run("success - Users found", func() {
		permissionTypes := []enum.PermissionType{enum.PermissionAll, enum.PermissionGeneral}
		s.userRepository.On("FindUsersByPermission", s.db, permissionTypes).Return([]*entity.User{}, true).Once()

		s.userService.GetUsersByPermission(permissionTypes)

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUsersByStatus() {
	s.Run("success - Users found", func() {
		request := userdto.GetUsersListRequest{
			Statuses: []uint{1, 2},
		}
		statuses := make([]enum.UserStatus, len(request.Statuses))
		for i, status := range request.Statuses {
			statuses[i] = enum.UserStatus(status)
		}
		users := []*entity.User{
			{
				Status:         enum.UserStatusActive,
				ProfilePicPath: "profile.jpg",
			},
			{
				Status:         enum.UserStatusBlock,
				ProfilePicPath: "profile.jpg",
			},
		}
		s.userRepository.On("FindUserByStatus", s.db, statuses, mock.Anything).Return(users, true).Once()
		s.s3Storage.On("GetPresignedURL", enum.ProfilePic, "profile.jpg", 8*time.Hour).Return("https://example.com/profile.jpg").Twice()

		s.userService.GetUsersByStatus(request)

		s.userRepository.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestBanUser() {
	s.Run("success - User banned", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		s.userService.BanUser(userID)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - User not found", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.BanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - User already banned", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{Status: enum.UserStatusBlock}, true).Once()

		s.Panics(func() {
			s.userService.BanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Update User Error", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(errors.New("update error")).Once()

		s.Panics(func() {
			s.userService.BanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUnbanUser() {
	s.Run("success - User unbanned", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		s.userService.UnbanUser(userID)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - User not found", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.UnbanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - User already unbanned", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{Status: enum.UserStatusActive}, true).Once()

		s.Panics(func() {
			s.userService.UnbanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Update User Error", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(errors.New("update error")).Once()

		s.Panics(func() {
			s.userService.UnbanUser(userID)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestRegister() {
	s.Run("success - User registered", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		otp := "123456"
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()
		s.userRepository.On("DeleteUserByPhone", s.db, mock.Anything).Return(nil).Once()
		s.userRepository.On("CreateUser", s.db, mock.Anything).Return(nil).Once()
		s.otpService.On("GenerateOTP").Return(otp, 1234567890).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, mock.Anything).Return(nil).Once()
		s.rabbitMQ.On("PublishMessage", mock.Anything, mock.Anything).Return(nil).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}
		s.userService.Register(request)

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.rabbitMQ.AssertExpectations(s.T())
	})
	s.Run("Error - duplicate phone number", func() {
		otpData := &userdto.OTPData{}
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}
		s.Panics(func() {
			s.userService.Register(request)
		})

		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Password too weak", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "weakpassword",
		}

		s.Panics(func() {
			s.userService.Register(request)
		})
		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Hash Password Error", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  strings.Repeat("A1@j", 100),
		}

		s.Panics(func() {
			s.userService.Register(request)
		})
		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Delete User Error", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		s.userRepository.On("DeleteUserByPhone", s.db, mock.Anything).Return(errors.New("delete error")).Once()
		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}
		s.Panics(func() {
			s.userService.Register(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Create User Error", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()
		s.userRepository.On("DeleteUserByPhone", s.db, mock.Anything).Return(nil).Once()

		s.userRepository.On("CreateUser", s.db, mock.Anything).Return(errors.New("create error")).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}
		s.Panics(func() {
			s.userService.Register(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Set OTP to Cache Error", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		otp := "123456"
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()
		s.userRepository.On("DeleteUserByPhone", s.db, mock.Anything).Return(nil).Once()
		s.userRepository.On("CreateUser", s.db, mock.Anything).Return(nil).Once()
		s.otpService.On("GenerateOTP").Return(otp, 1234567890).Once()

		s.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, mock.Anything).Return(errors.New("cache error")).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}

		s.Panics(func() {
			s.userService.Register(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Publish Message Error", func() {
		var nilOTPData *userdto.OTPData = nil
		var nilUser *entity.User = nil
		otp := "123456"
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()
		s.userRepository.On("DeleteUserByPhone", s.db, mock.Anything).Return(nil).Once()
		s.userRepository.On("CreateUser", s.db, mock.Anything).Return(nil).Once()
		s.otpService.On("GenerateOTP").Return(otp, 1234567890).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, mock.Anything).Return(nil).Once()
		s.rabbitMQ.On("PublishMessage", mock.Anything, mock.Anything).Return(errors.New("publish message error")).Once()

		request := userdto.BasicRegisterRequest{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Password:  "Password@123",
		}
		s.Panics(func() {
			s.userService.Register(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.rabbitMQ.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyPhone() {
	s.Run("success - Phone verified", func() {
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(&entity.User{}, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(nil).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   "123456",
		}

		s.userService.VerifyPhone(request)

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("Error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   "123456",
		}

		s.Panics(func() {
			s.userService.VerifyPhone(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - Phone already verified", func() {
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(&entity.User{PhoneVerified: true}, true).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   "123456",
		}

		s.Panics(func() {
			s.userService.VerifyPhone(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("Error - OTP verification failed", func() {
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(&entity.User{}, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(errors.New("invalid OTP")).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   "123456",
		}

		s.Panics(func() {
			s.userService.VerifyPhone(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("Error - Update User Error", func() {
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(&entity.User{}, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, mock.Anything).Return(errors.New("update error")).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   "123456",
		}

		s.Panics(func() {
			s.userService.VerifyPhone(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestFindUserPermissions() {
	s.Run("success - User permissions found", func() {
		user := &entity.User{
			Roles: []entity.Role{
				{
					Name: "admin",
				},
				{
					Name: "common",
				},
			},
		}
		s.userRepository.On("FindUserRoles", s.db, user).Return(nil)
		for _, role := range user.Roles {
			s.userRepository.On("FindRolePermissions", s.db, &role).Return(nil)
		}
		s.userService.FindUserPermissions(user)

		s.userRepository.AssertExpectations(s.T())

	})
	s.Run("Error - User roles not found", func() {
		user := &entity.User{
			Roles: []entity.Role{},
		}
		s.userRepository.On("FindUserRoles", s.db, user).Return(errors.New("roles not found"))

		s.Panics(func() {
			s.userService.FindUserPermissions(user)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestLogin() {
	s.Run("success - User logged in", func() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), 14)
		mockAccessToken := "mock-access-token"
		mockRefreshToken := "mock-refresh-token"
		user := &entity.User{
			FirstName:     "John",
			LastName:      "Doe",
			PhoneVerified: true,
			Password:      string(hashedPassword),
			Roles: []entity.Role{
				{
					Name: "admin",
				},
				{
					Name: "common",
				},
			},
		}

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()
		s.jwtService.On("GenerateToken", mock.Anything).Return(mockAccessToken, mockRefreshToken).Once()
		s.userRepository.On("FindUserRoles", s.db, user).Return(nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, mock.Anything).Return(nil).Twice()

		request := userdto.LoginRequest{
			Phone:    "1234567890",
			Password: "Password@123",
		}
		response, _ := s.userService.Login(request)

		s.Equal(response.AccessToken, mockAccessToken)
		s.Equal(response.RefreshToken, mockRefreshToken)
		s.Equal(response.FirstName, user.FirstName)
		s.Equal(response.LastName, user.LastName)

		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
	s.Run("error - Wrong Password", func() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), 14)
		user := &entity.User{
			PhoneVerified: true,
			Password:      string(hashedPassword),
		}

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.LoginRequest{
			Phone:    "1234567890",
			Password: "Password@1234",
		}
		s.Panics(func() {
			s.userService.Login(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.LoginRequest{
			Phone:    "1234567890",
			Password: "Password@1234",
		}
		s.Panics(func() {
			s.userService.Login(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), 14)
		user := &entity.User{
			PhoneVerified: false,
			Password:      string(hashedPassword),
		}

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.LoginRequest{
			Phone:    "1234567890",
			Password: "Password@1234",
		}
		s.Panics(func() {
			s.userService.Login(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestForgotPassword() {
	s.Run("success - OTP sent", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		otp := "123456"
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("GenerateOTP").Return(otp, 2)
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, mock.Anything).Return(nil).Once()

		request := userdto.ForgotPasswordRequest{
			Phone: "1234567890",
		}
		s.userService.ForgotPassword(request)

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.ForgotPasswordRequest{
			Phone: "1234567890",
		}
		s.Panics(func() {
			s.userService.ForgotPassword(request)
		})
		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		user := &entity.User{
			PhoneVerified: false,
		}
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.ForgotPasswordRequest{
			Phone: "1234567890",
		}
		s.Panics(func() {
			s.userService.ForgotPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("error - Set OTP to cache error", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		otp := "123456"
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("GenerateOTP").Return(otp, 2)
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, otp, mock.Anything).Return(errors.New("test error")).Once()

		request := userdto.ForgotPasswordRequest{
			Phone: "1234567890",
		}
		s.Panics(func() {
			s.userService.ForgotPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyOTP() {
	s.Run("success - OTP verified", func() {
		user := &entity.User{
			FirstName:     "John",
			LastName:      "Doe",
			PhoneVerified: true,
			Roles: []entity.Role{
				{
					Name: "admin",
				},
				{
					Name: "common",
				},
			},
		}
		mockAccessToken := "mock-access-token"
		mockRefreshToken := "mock-refresh-token"
		otp := "123456"

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, otp).Return(nil).Once()
		s.jwtService.On("GenerateToken", mock.Anything).Return(mockAccessToken, mockRefreshToken).Once()
		s.userRepository.On("FindUserRoles", s.db, user).Return(nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, mock.Anything).Return(nil).Twice()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   otp,
		}
		response, _ := s.userService.VerifyOTP(request)

		s.Equal(response.AccessToken, mockAccessToken)
		s.Equal(response.RefreshToken, mockRefreshToken)
		s.Equal(response.FirstName, user.FirstName)
		s.Equal(response.LastName, user.LastName)

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		otp := "123456"
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   otp,
		}
		s.Panics(func() {
			s.userService.VerifyOTP(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.jwtService.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		user := &entity.User{
			PhoneVerified: false,
		}
		otp := "123456"
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   otp,
		}
		s.Panics(func() {
			s.userService.VerifyOTP(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("error - OTP verification failed", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		otp := "123456"
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, otp).Return(errors.New("invalid OTP")).Once()

		request := userdto.VerifyPhoneRequest{
			Phone: "1234567890",
			OTP:   otp,
		}
		s.Panics(func() {
			s.userService.VerifyOTP(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestCompleteRegister() {
	s.Run("success - User registered", func() {
		user := &entity.User{}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()

		request := userdto.CompleteRegisterRequest{
			UserID:       1,
			NationalCode: "1234567890",
			ProfilePic:   nil,
			TemplateFile: "template.html",
			EmailSubject: "Welcome",
		}
		s.userService.CompleteRegister(request)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("success - User entered new email", func() {
		user := &entity.User{}
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(&userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()
		s.otpService.On("GenerateOTP").Return("123456", 2).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		request := userdto.CompleteRegisterRequest{
			UserID: 1,
			Email:  "test@example.com",
		}
		s.userService.CompleteRegister(request)

		s.userRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
	})
	s.Run("success - User entered new profile pic", func() {
		user := &entity.User{}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.ProfilePic, mock.Anything, mock.Anything).Return().Once()

		request := userdto.CompleteRegisterRequest{
			UserID: 1,
			ProfilePic: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     int64(len([]byte("test"))),
			},
		}
		s.userService.CompleteRegister(request)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.CompleteRegisterRequest{
			UserID: 1,
		}
		s.Panics(func() {
			s.userService.CompleteRegister(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Update User Error", func() {
		user := &entity.User{}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(errors.New("update error")).Once()

		request := userdto.CompleteRegisterRequest{
			UserID: 1,
		}
		s.Panics(func() {
			s.userService.CompleteRegister(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestVerifyEmail() {
	s.Run("success - Email verified", func() {
		user := &entity.User{
			PhoneVerified: true,
			EmailVerified: false,
		}
		otp := "123456"

		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    otp,
		}

		s.userService.VerifyEmail(request)

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    "123456",
		}
		s.Panics(func() {
			s.userService.VerifyEmail(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		user := &entity.User{
			PhoneVerified: false,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    "123456",
		}
		s.Panics(func() {
			s.userService.VerifyEmail(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("error - Email already verified", func() {
		user := &entity.User{
			PhoneVerified: true,
			EmailVerified: true,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    "123456",
		}
		s.Panics(func() {
			s.userService.VerifyEmail(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - OTP verification failed", func() {
		user := &entity.User{
			PhoneVerified: true,
			EmailVerified: false,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(errors.New("invalid OTP")).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    "123456",
		}
		s.Panics(func() {
			s.userService.VerifyEmail(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
	s.Run("error - Update User Error", func() {
		user := &entity.User{
			PhoneVerified: true,
			EmailVerified: false,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.otpService.On("VerifyOTP", mock.Anything, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(errors.New("update error")).Once()

		request := userdto.VerifyEmailRequest{
			UserID: 1,
			Email:  "test@example.com",
			OTP:    "123456",
		}
		s.Panics(func() {
			s.userService.VerifyEmail(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.otpService.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestResetPassword() {
	s.Run("success - Password reset", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword@123",
		}
		s.userService.ResetPassword(request)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword@123",
		}
		s.Panics(func() {
			s.userService.ResetPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		user := &entity.User{
			PhoneVerified: false,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword@123",
		}
		s.Panics(func() {
			s.userService.ResetPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Update User Error", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(errors.New("update error")).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "NewPassword@123",
		}
		s.Panics(func() {
			s.userService.ResetPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Hash Password Error", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: strings.Repeat("A1@j", 100),
		}
		s.Panics(func() {
			s.userService.ResetPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Password too weak", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()

		request := userdto.ResetPasswordRequest{
			UserID:   1,
			Password: "weakpassword",
		}
		s.Panics(func() {
			s.userService.ResetPassword(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestFindUserByPhone() {
	s.Run("success - User found", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		phone := "1234567890"

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		s.userService.FindActiveUserByPhone(phone)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		phone := "1234567890"

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.FindActiveUserByPhone(phone)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Phone not verified", func() {
		user := &entity.User{
			PhoneVerified: false,
		}
		phone := "1234567890"

		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(user, true).Once()

		s.Panics(func() {
			s.userService.FindActiveUserByPhone(phone)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUpdateProfile() {
	s.Run("success - Profile updated", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(&userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()
		s.otpService.On("GenerateOTP").Return("123456", 2).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.ProfilePic, mock.Anything, mock.Anything).Return().Once()
		s.s3Storage.On("DeleteObject", enum.ProfilePic, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(nil).Once()

		stringPtr := func(s string) *string {
			return &s
		}

		request := userdto.UpdateProfileRequest{
			UserID:       1,
			FirstName:    stringPtr("John"),
			LastName:     stringPtr("Doe"),
			Email:        stringPtr("test@example.com"),
			NationalCode: stringPtr("1234567890"),
			ProfilePic: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     int64(len([]byte("test"))),
			},
			TemplateFile: "template.html",
			EmailSubject: "Welcome",
		}
		s.userService.UpdateProfile(request)

		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.UpdateProfileRequest{
			UserID: 1,
		}
		s.Panics(func() {
			s.userService.UpdateProfile(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - S3 error", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(&userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()
		s.otpService.On("GenerateOTP").Return("123456", 2).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.ProfilePic, mock.Anything, mock.Anything).Return().Once()
		s.s3Storage.On("DeleteObject", enum.ProfilePic, mock.Anything).Return(errors.New("S3 error")).Once()

		stringPtr := func(s string) *string {
			return &s
		}

		request := userdto.UpdateProfileRequest{
			UserID:       1,
			FirstName:    stringPtr("John"),
			LastName:     stringPtr("Doe"),
			Email:        stringPtr("test@example.com"),
			NationalCode: stringPtr("1234567890"),
			ProfilePic: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     int64(len([]byte("test"))),
			},
			TemplateFile: "template.html",
			EmailSubject: "Welcome",
		}
		s.Panics(func() {
			s.userService.UpdateProfile(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
	s.Run("error - Update User Error", func() {
		user := &entity.User{
			PhoneVerified: true,
		}
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(&userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}, false).Once()
		s.userRepository.On("FindUserByEmail", s.db, mock.Anything).Return(nilUser, false).Once()
		s.otpService.On("GenerateOTP").Return("123456", 2).Once()
		s.userCacheRepository.On("Set", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.emailService.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.ProfilePic, mock.Anything, mock.Anything).Return().Once()
		s.s3Storage.On("DeleteObject", enum.ProfilePic, mock.Anything).Return(nil).Once()
		s.userRepository.On("UpdateUser", s.db, user).Return(errors.New("update error")).Once()

		stringPtr := func(s string) *string {
			return &s
		}

		request := userdto.UpdateProfileRequest{
			UserID:       1,
			FirstName:    stringPtr("John"),
			LastName:     stringPtr("Doe"),
			Email:        stringPtr("test@example.com"),
			NationalCode: stringPtr("1234567890"),
			ProfilePic: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     int64(len([]byte("test"))),
			},
			TemplateFile: "template.html",
			EmailSubject: "Welcome",
		}
		s.Panics(func() {
			s.userService.UpdateProfile(request)
		})

		s.userRepository.AssertExpectations(s.T())
		s.userCacheRepository.AssertExpectations(s.T())
		s.emailService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})

}

func (s *UserServiceTestSuite) TestGetAllPermissions() {
	s.Run("success - Permissions found", func() {
		permissions := []*entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
			{
				Type:        enum.PermissionAll,
				Description: "دسترسی کامل به سیستم",
				Category:    enum.CategoryGeneral,
			},
		}
		s.userRepository.On("FindAllPermissions", s.db).Return(permissions, nil).Once()
		response, _ := s.userService.GetAllPermissions()

		s.Equal(response[0].Name, enum.PermissionGeneral.String())
		s.Equal(response[1].Name, enum.PermissionAll.String())
		s.Equal(response[0].Description, "دسترسی عمومی")
		s.Equal(response[1].Description, "دسترسی کامل به سیستم")
		s.Equal(response[0].Category, enum.CategoryGeneral.String())
		s.Equal(response[1].Category, enum.CategoryGeneral.String())

		s.userRepository.AssertExpectations(s.T())
	})

}

func (s *UserServiceTestSuite) TestGetRolePermissions() {
	s.Run("success - Role permissions found", func() {
		role := &entity.Role{
			Name: "admin",
		}
		role.Permissions = []entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
			{
				Type:        enum.PermissionAll,
				Description: "دسترسی کامل به سیستم",
				Category:    enum.CategoryGeneral,
			},
		}
		s.userRepository.On("FindRolePermissions", s.db, role).Return(nil).Once()

		response, _ := s.userService.getRolePermissions(role)

		s.Equal(response[0].Name, enum.PermissionGeneral.String())
		s.Equal(response[1].Name, enum.PermissionAll.String())
		s.Equal(response[0].Description, "دسترسی عمومی")
		s.Equal(response[1].Description, "دسترسی کامل به سیستم")
		s.Equal(response[0].Category, enum.CategoryGeneral.String())
		s.Equal(response[1].Category, enum.CategoryGeneral.String())

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - repo error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		s.userRepository.On("FindRolePermissions", s.db, role).Return(errors.New("repo error")).Once()
		s.Panics(func() {
			s.userService.getRolePermissions(role)
		})
		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetAllRoles() {
	s.Run("success - Roles found", func() {
		roles := []*entity.Role{
			{
				Name: "admin",
			},
			{
				Name: "user",
			},
		}
		roles[0].Permissions = []entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}
		roles[1].Permissions = []entity.Permission{
			{
				Type:        enum.PermissionAll,
				Description: "دسترسی کامل به سیستم",
				Category:    enum.CategoryGeneral,
			},
		}
		s.userRepository.On("FindAllRoles", s.db).Return(roles, nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, roles[0]).Return(nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, roles[1]).Return(nil).Once()

		response, _ := s.userService.GetAllRoles()

		s.Equal(response[0].Name, "admin")
		s.Equal(response[1].Name, "user")
		s.Equal(response[0].Permissions[0].Name, enum.PermissionGeneral.String())
		s.Equal(response[1].Permissions[0].Name, enum.PermissionAll.String())
		s.Equal(response[0].Permissions[0].Description, "دسترسی عمومی")
		s.Equal(response[1].Permissions[0].Description, "دسترسی کامل به سیستم")
		s.Equal(response[0].Permissions[0].Category, enum.CategoryGeneral.String())
		s.Equal(response[1].Permissions[0].Category, enum.CategoryGeneral.String())

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestCreateRole() {
	s.Run("success - Role created", func() {
		role := &entity.Role{
			Name: "admin",
		}
		var nilRole *entity.Role = nil
		permissions := []*entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
			{
				Type:        enum.PermissionAll,
				Description: "دسترسی کامل به سیستم",
				Category:    enum.CategoryGeneral,
			},
		}
		role.Permissions = nil
		s.userRepository.On("FindRoleByName", s.db, "admin").Return(nilRole, false).Once()
		s.userRepository.On("CreateRole", s.db, role).Return(nil).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(permissions[0], true).Once()
		s.userRepository.On("AssignPermissionToRole", s.db, role, permissions[0]).Return(nil).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(permissions[1], true).Once()
		s.userRepository.On("AssignPermissionToRole", s.db, role, permissions[1]).Return(nil).Once()

		s.userService.CreateRole(userdto.NewRoleRequest{
			Name:          "admin",
			PermissionIDs: []uint{1, 2, 1},
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Role already exists", func() {
		role := &entity.Role{
			Name: "admin",
		}

		s.userRepository.On("FindRoleByName", s.db, "admin").Return(role, true).Once()

		request := userdto.NewRoleRequest{
			Name:          "admin",
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.CreateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Create Role Error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		s.userRepository.On("FindRoleByName", s.db, "admin").Return(role, false).Once()
		s.userRepository.On("CreateRole", s.db, role).Return(errors.New("create role error")).Once()

		request := userdto.NewRoleRequest{
			Name:          "admin",
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.CreateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Permission not found", func() {
		role := &entity.Role{
			Name: "admin",
		}
		var nilPermission *entity.Permission = nil
		s.userRepository.On("FindRoleByName", s.db, "admin").Return(role, false).Once()
		s.userRepository.On("CreateRole", s.db, role).Return(nil).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(nilPermission, false).Once()

		request := userdto.NewRoleRequest{
			Name:          "admin",
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.CreateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Assign Permission to Role Error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		permission := &entity.Permission{
			Type:        enum.PermissionGeneral,
			Description: "دسترسی عمومی",
			Category:    enum.CategoryGeneral,
		}
		s.userRepository.On("FindRoleByName", s.db, "admin").Return(role, false).Once()
		s.userRepository.On("CreateRole", s.db, role).Return(nil).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(permission, true).Once()
		s.userRepository.On("AssignPermissionToRole", s.db, role, permission).Return(errors.New("assign permission to role error")).Once()

		request := userdto.NewRoleRequest{
			Name:          "admin",
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.CreateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetRoleDetails() {
	s.Run("success - Room details found", func() {
		role := &entity.Role{
			Name: "admin",
		}
		role.Permissions = []entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("FindRolePermissions", s.db, role).Return(nil).Once()
		response, _ := s.userService.GetRoleDetails(1)

		s.Equal(response.Name, "admin")
		s.Equal(response.Permissions[0].Name, enum.PermissionGeneral.String())
	})
	s.Run("error - Role not found", func() {
		var nilRole *entity.Role = nil
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(nilRole, false).Once()

		s.Panics(func() {
			s.userService.GetRoleDetails(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetRoleOwners() {
	s.Run("success - Role owners found", func() {
		role := &entity.Role{
			Name: "admin",
		}
		users := []*entity.User{
			{
				Phone:          "09123456789",
				ProfilePicPath: "profile.jpg",
			},
			{
				Phone:          "09123456789",
				ProfilePicPath: "profile.jpg",
			},
		}
		role.Permissions = []entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}

		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("FindUsersByRoleID", s.db, mock.Anything).Return(users, nil).Once()
		s.s3Storage.On("GetPresignedURL", enum.ProfilePic, mock.Anything, mock.Anything).Return("https://example.com/profile.jpg").Twice()

		response, _ := s.userService.GetRoleOwners(1)

		s.Equal(response[0].Phone, "09123456789")
		s.Equal(response[1].Phone, "09123456789")
	})
	s.Run("error - Role not found", func() {
		var nilRole *entity.Role = nil
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(nilRole, false).Once()

		s.Panics(func() {
			s.userService.GetRoleOwners(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestGetUserRoles() {
	s.Run("success - User roles found", func() {
		user := &entity.User{
			Roles: []entity.Role{
				{
					Name: "admin",
				},
				{
					Name: "user",
				},
			},
		}
		user.Roles[0].Permissions = []entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}
		user.Roles[1].Permissions = []entity.Permission{
			{
				Type:        enum.PermissionAll,
				Description: "دسترسی کامل به سیستم",
				Category:    enum.CategoryGeneral,
			},
		}

		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("FindUserRoles", s.db, user).Return(nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, &user.Roles[0]).Return(nil).Once()
		s.userRepository.On("FindRolePermissions", s.db, &user.Roles[1]).Return(nil).Once()

		response, _ := s.userService.GetUserRoles(1)

		s.Equal(response[0].Name, "admin")
		s.Equal(response[1].Name, "user")
		s.Equal(response[0].Permissions[0].Name, enum.PermissionGeneral.String())
		s.Equal(response[1].Permissions[0].Name, enum.PermissionAll.String())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.GetUserRoles(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Find User Roles Error", func() {
		user := &entity.User{}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("FindUserRoles", s.db, user).Return(errors.New("find user roles error")).Once()

		s.Panics(func() {
			s.userService.GetUserRoles(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestDeleteRole() {
	s.Run("success - Role deleted", func() {
		role := &entity.Role{
			Name: "admin",
		}
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("DeleteRole", s.db, mock.Anything).Return(nil).Once()

		s.userService.DeleteRole(1)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Role not found", func() {
		var nilRole *entity.Role = nil
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(nilRole, false).Once()

		s.Panics(func() {
			s.userService.DeleteRole(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Delete Role Error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("DeleteRole", s.db, mock.Anything).Return(errors.New("delete role error")).Once()

		s.Panics(func() {
			s.userService.DeleteRole(1)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUpdateRole() {
	s.Run("success - Role updated", func() {
		role := &entity.Role{
			Name: "admin",
		}
		permissions := []*entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}

		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("ReplaceRolePermissions", s.db, role, mock.Anything).Return(nil).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(permissions[0], true).Once()

		request := userdto.UpdateRoleRequest{
			RoleID:        1,
			PermissionIDs: []uint{1, 1},
		}
		s.userService.UpdateRole(request)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Role not found", func() {
		var nilRole *entity.Role = nil
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(nilRole, false).Once()

		s.Panics(func() {
			s.userService.UpdateRole(userdto.UpdateRoleRequest{RoleID: 1, PermissionIDs: []uint{1, 1}})
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Role name update error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		newName := "admin2"
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("UpdateRole", s.db, role).Return(errors.New("update role error")).Once()

		request := userdto.UpdateRoleRequest{
			RoleID: 1,
			Name:   &newName,
		}
		s.Panics(func() {
			s.userService.UpdateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Permission not found", func() {
		role := &entity.Role{
			Name: "admin",
		}
		var nilPermission *entity.Permission = nil

		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(nilPermission, false).Once()

		request := userdto.UpdateRoleRequest{
			RoleID:        1,
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.UpdateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Assign Permission to Role Error", func() {
		role := &entity.Role{
			Name: "admin",
		}
		permissions := []*entity.Permission{
			{
				Type:        enum.PermissionGeneral,
				Description: "دسترسی عمومی",
				Category:    enum.CategoryGeneral,
			},
		}
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("ReplaceRolePermissions", s.db, role, mock.Anything).Return(errors.New("assign permission to role error")).Once()
		s.userRepository.On("FindPermissionByID", s.db, mock.Anything).Return(permissions[0], true).Once()

		request := userdto.UpdateRoleRequest{
			RoleID:        1,
			PermissionIDs: []uint{1},
		}
		s.Panics(func() {
			s.userService.UpdateRole(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func (s *UserServiceTestSuite) TestUpdateUserRoles() {
	s.Run("success - User roles updated", func() {
		user := &entity.User{
			Roles: []entity.Role{
				{
					Name: "admin",
				},
				{
					Name: "user",
				},
			},
		}
		role := &entity.Role{
			Name: "admin2",
		}

		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("ReplaceUserRoles", s.db, user, mock.Anything).Return(nil).Once()

		request := userdto.UpdateUserRolesRequest{
			UserID:  1,
			RoleIDs: []uint{1, 1},
		}
		s.userService.UpdateUserRoles(request)

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not found", func() {
		var nilUser *entity.User = nil
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(nilUser, false).Once()

		request := userdto.UpdateUserRolesRequest{
			UserID:  1,
			RoleIDs: []uint{1, 1},
		}
		s.Panics(func() {
			s.userService.UpdateUserRoles(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Role not found", func() {
		user := &entity.User{}
		var nilRole *entity.Role = nil

		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(nilRole, false).Once()

		request := userdto.UpdateUserRolesRequest{
			UserID:  1,
			RoleIDs: []uint{1, 1},
		}
		s.Panics(func() {
			s.userService.UpdateUserRoles(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
	s.Run("error - Assign Role to User Error", func() {
		user := &entity.User{}
		role := &entity.Role{
			Name: "admin",
		}
		s.userRepository.On("FindUserByID", s.db, mock.Anything).Return(user, true).Once()
		s.userRepository.On("FindRoleByID", s.db, mock.Anything).Return(role, true).Once()
		s.userRepository.On("ReplaceUserRoles", s.db, user, mock.Anything).Return(errors.New("assign role to user error")).Once()

		request := userdto.UpdateUserRolesRequest{
			UserID:  1,
			RoleIDs: []uint{1, 1},
		}
		s.Panics(func() {
			s.userService.UpdateUserRoles(request)
		})

		s.userRepository.AssertExpectations(s.T())
	})
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
