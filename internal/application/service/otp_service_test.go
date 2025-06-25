package service

import (
	"context"
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	mocks "github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OtpServiceTestSuite struct {
	suite.Suite
	constants           *bootstrap.Constants
	otpConfig           *bootstrap.OTP
	userCacheRepository *mocks.UserCacheRepositoryMock
	otpService          *OTPService
}

func (s *OtpServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	s.constants = config.Constants
	s.otpConfig = &bootstrap.OTP{
		Length:       6,
		ExpiryMinute: 10,
		MaxAttempts:  3,
	}
	s.userCacheRepository = mocks.NewUserCacheRepositoryMock()
	s.otpService = NewOTPService(s.constants, s.otpConfig, s.userCacheRepository)
}

func (s *OtpServiceTestSuite) TestGenerateOTP() {
	s.Run("success - OTP generated", func() {
		otp, _, _ := s.otpService.GenerateOTP()
		s.Equal(len(otp), s.otpConfig.Length)
	})
}

func (s *OtpServiceTestSuite) TestVerifyOTP() {
	s.Run("success - OTP verified", func() {
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()
		err := s.otpService.VerifyOTP(mock.Anything, "123456")
		s.NoError(err)
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("error - OTP not found", func() {
		var nilOTPData *userdto.OTPData = nil

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(nilOTPData, false).Once()
		err := s.otpService.VerifyOTP(mock.Anything, "123456")

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
	})
	s.Run("error - OTP is invalid", func() {
		otpData := &userdto.OTPData{
			OTP:      "123456",
			Attempts: 0,
		}

		s.userCacheRepository.On("Get", context.Background(), mock.Anything).Return(otpData, true).Once()
		err := s.otpService.VerifyOTP(mock.Anything, "123457")

		s.Error(err)
		s.userCacheRepository.AssertExpectations(s.T())
	})
}

func TestOtpService(t *testing.T) {
	suite.Run(t, new(OtpServiceTestSuite))
}
