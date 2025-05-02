package serviceimpl_test

import (
	"github.com/BargheNo/Backend/bootstrap"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/mocks"
)

type UserServiceTestSuite struct {
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
