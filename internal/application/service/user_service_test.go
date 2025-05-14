package serviceimpl

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
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
	}
	s.userService = NewUserService(deps)
}

func (s *UserServiceTestSuite) TestDoesUserExist() {
	s.Run("success - User exists", func() {
		userID := uint(1)
		s.userRepository.On("FindUserByID", s.db, userID).Return(&entity.User{}, true).Once()

		s.userService.DoesUserExist(userID)

		s.userRepository.AssertExpectations(s.T())
	})

	s.Run("Error - User does not exist", func() {
		userID := uint(1)
		var nilUser *entity.User = nil

		s.userRepository.On("FindUserByID", s.db, userID).Return(nilUser, false).Once()

		s.Panics(func() {
			s.userService.DoesUserExist(userID)
		})

		s.userRepository.AssertExpectations(s.T())
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
	})
	s.Run("Error - duplicate phone number(pending for registration)", func() {
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
	s.Run("Error - duplicate phone number(verified)", func() {
		var nilOTPData *userdto.OTPData = nil
		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		s.userRepository.On("FindUserByPhone", s.db, mock.Anything).Return(&entity.User{PhoneVerified: true}, true).Once()

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
		response := s.userService.Login(request)

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
func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
