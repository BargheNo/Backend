package serviceimpl

import (
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CorporationServiceTestSuite struct {
	suite.Suite
	constants             *bootstrap.Constants
	userService           *mocks.UserServiceMock
	addressService        *mocks.AddressServiceMock
	s3Storage             *mocks.S3StorageMock
	corporationRepository *mocks.CorporationRepositoryMock
	db                    *mocks.DatabaseMock
	corporationService    *CorporationService
}

func (s *CorporationServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	s.constants = config.Constants
	s.userService = mocks.NewUserServiceMock()
	s.addressService = mocks.NewAddressServiceMock()
	s.s3Storage = mocks.NewS3StorageMock()
	s.corporationRepository = mocks.NewCorporationRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.corporationService = NewCorporationService(
		s.constants,
		s.userService,
		s.addressService,
		s.s3Storage,
		s.corporationRepository,
		s.db,
	)
}

func (s *CorporationServiceTestSuite) TestGetCorporationByIDAndStatus() {
	s.Run("success - Corporation found", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()

		response := s.corporationService.getCorporationByIDAndStatus(uint(1), enum.CorpStatusApproved)

		s.Equal(response.Status, enum.CorpStatusApproved)
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation not found", func() {
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(nilCorporation, false).Once()

		s.Panics(func() {
			s.corporationService.getCorporationByIDAndStatus(uint(1), enum.CorpStatusApproved)
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation status not approved", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()

		s.Panics(func() {
			s.corporationService.getCorporationByIDAndStatus(uint(1), enum.CorpStatusApproved)
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestDoesCorporationExist() {
	s.Run("success - Corporation found", func() {
		corporation := &entity.Corporation{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()

		s.corporationService.DoesCorporationExist(uint(1))

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation not found", func() {
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(nilCorporation, false).Once()

		s.Panics(func() {
			s.corporationService.DoesCorporationExist(uint(1))
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationCredentials() {
	s.Run("success - Corporation found", func() {
		corporation := &entity.Corporation{
			Name: "testName",
			// Logo:   "testLogo",
			Status: enum.CorpStatusApproved,
		}
		corporation.Addresses = []entity.Address{
			{
				PostalCode: "testPostalCode",
			},
		}
		corporation.ContactInformation = []entity.ContactInformation{
			{
				Value: "testValue",
			},
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.addressService.On("GetAddresses", mock.Anything).Return([]addressdto.AddressResponse{
			{
				PostalCode: "testPostalCode",
			},
		}).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return([]*entity.ContactInformation{
			{
				Value: "testValue",
			},
		}).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, mock.Anything).Return(&entity.ContactType{}, true)
		response := s.corporationService.GetCorporationCredentials(uint(1))

		s.Equal(response.ID, corporation.ID)
		s.Equal(response.Name, corporation.Name)
		// s.Equal(response.Logo, corporation.Logo)
		s.Equal(response.ContactInfo[0].Value, corporation.ContactInformation[0].Value)
		s.Equal(response.Addresses[0].PostalCode, corporation.Addresses[0].PostalCode)

		s.addressService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation not found", func() {
		var nilCorporation *entity.Corporation = nil
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(nilCorporation, false).Once()

		s.Panics(func() {
			s.corporationService.GetCorporationCredentials(uint(1))
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestISCorporationApproved() {
	s.Run("success - Corporation is approved", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()

		response := s.corporationService.ISCorporationApproved(uint(1))

		s.True(response)

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation not found", func() {
		var nilCorporation *entity.Corporation = nil
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(nilCorporation, false).Once()

		s.Panics(func() {
			s.corporationService.ISCorporationApproved(uint(1))
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestCheckApplicantAccess() {
	s.Run("success - Applicant has access", func() {
		staff := &entity.CorporationStaff{}
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(staff, true).Once()

		s.corporationService.CheckApplicantAccess(uint(1), uint(1))

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Applicant does not have access", func() {
		var nilStaff *entity.CorporationStaff = nil
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(nilStaff, false).Once()

		s.Panics(func() {
			s.corporationService.CheckApplicantAccess(uint(1), uint(1))
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func TestCorporationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CorporationServiceTestSuite))
}
