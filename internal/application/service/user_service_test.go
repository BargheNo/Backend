package serviceimpl_test

import (
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/suite"
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
	userService         *serviceimpl.UserService
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

	deps := serviceimpl.UserServiceDeps{
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
	s.userService = serviceimpl.NewUserService(deps)
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

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
