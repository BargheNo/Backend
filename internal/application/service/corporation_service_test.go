package service

import (
	"errors"
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"

	databaseMocks "github.com/BargheNo/Backend/mocks/infrastructure/database"
	repositoryMocks "github.com/BargheNo/Backend/mocks/infrastructure/repository/postgres"
	s3Mocks "github.com/BargheNo/Backend/mocks/infrastructure/s3"
	serviceMocks "github.com/BargheNo/Backend/mocks/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CorporationServiceTestSuite struct {
	suite.Suite
	constants             *bootstrap.Constants
	userService           *serviceMocks.UserServiceMock
	addressService        *serviceMocks.AddressServiceMock
	rbacService           *serviceMocks.RBACServiceMock
	s3Storage             *s3Mocks.S3StorageMock
	corporationRepository *repositoryMocks.CorporationRepositoryMock
	db                    *databaseMocks.DatabaseMock
	corporationService    *CorporationService
}

func (suite *CorporationServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	suite.constants = config.Constants
	suite.userService = serviceMocks.NewUserServiceMock()
	suite.addressService = serviceMocks.NewAddressServiceMock()
	suite.rbacService = serviceMocks.NewRBACServiceMock()
	suite.s3Storage = s3Mocks.NewS3StorageMock()
	suite.corporationRepository = repositoryMocks.NewCorporationRepositoryMock()
	suite.db = databaseMocks.NewDatabaseMock()

	corporationService := NewCorporationService(
		suite.constants,
		suite.userService,
		suite.addressService,
		suite.rbacService,
		suite.s3Storage,
		suite.corporationRepository,
		suite.db,
	)

	suite.corporationService = corporationService
}

func (s *CorporationServiceTestSuite) TestGetCorporationStatuses() {
	s.Run("success - returns corporation statuses", func() {
		result := s.corporationService.GetCorporationStatuses()

		s.NotNil(result)
		s.Greater(len(result), 0)

		// Verify each status has ID and Name
		for _, status := range result {
			s.Greater(status.ID, uint(0))
			s.NotEmpty(status.Name)
		}
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationSortableColumns() {
	s.Run("success - returns corporation sortable columns", func() {
		result := s.corporationService.GetCorporationSortableColumns()

		s.NotNil(result)
		s.Greater(len(result), 0)

		// Verify each column has ID and Name
		for _, col := range result {
			s.Greater(col.ID, uint(0))
			s.NotEmpty(col.Name)
		}
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationStaffSortableColumns() {
	s.Run("success - returns staff sortable columns", func() {
		result := s.corporationService.GetCorporationStaffSortableColumns()

		s.NotNil(result)
		s.Greater(len(result), 0)

		for _, col := range result {
			s.Greater(col.ID, uint(0))
			s.NotEmpty(col.Name)
		}
	})
}

func (s *CorporationServiceTestSuite) TestDoesCorporationExist() {
	s.Run("success - corporation exists", func() {
		corporationID := uint(1)
		corporation := &entity.Corporation{
			Name: "Test Corp",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(corporation, nil).Once()

		err := s.corporationService.DoesCorporationExist(corporationID)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		corporationID := uint(1)
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(nilCorporation, nil).Once()

		err := s.corporationService.DoesCorporationExist(corporationID)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		corporationID := uint(1)
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(nilCorporation, errors.New("database error")).Once()

		err := s.corporationService.DoesCorporationExist(corporationID)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestCheckApplicantAccess() {
	s.Run("success - applicant has access", func() {
		corporationID := uint(1)
		applicantID := uint(2)
		staff := &entity.CorporationStaff{
			UserID:        applicantID,
			CorporationID: corporationID,
		}

		s.corporationRepository.On("FindCorporationStaff", s.db, applicantID, corporationID).Return(staff, nil).Once()

		err := s.corporationService.CheckApplicantAccess(corporationID, applicantID)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - applicant not found", func() {
		corporationID := uint(1)
		applicantID := uint(2)
		var nilStaff *entity.CorporationStaff = nil

		s.corporationRepository.On("FindCorporationStaff", s.db, applicantID, corporationID).Return(nilStaff, nil).Once()

		err := s.corporationService.CheckApplicantAccess(corporationID, applicantID)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		corporationID := uint(1)
		applicantID := uint(2)
		var nilStaff *entity.CorporationStaff = nil

		s.corporationRepository.On("FindCorporationStaff", s.db, applicantID, corporationID).Return(nilStaff, errors.New("database error")).Once()

		err := s.corporationService.CheckApplicantAccess(corporationID, applicantID)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestISCorporationApproved() {
	s.Run("success - corporation is approved", func() {
		corporationID := uint(1)
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(corporation, nil).Once()

		err := s.corporationService.ISCorporationApproved(corporationID)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - corporation not approved", func() {
		corporationID := uint(1)
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(corporation, nil).Once()

		err := s.corporationService.ISCorporationApproved(corporationID)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		corporationID := uint(1)
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(nilCorporation, nil).Once()

		err := s.corporationService.ISCorporationApproved(corporationID)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationCredentials() {
	s.Run("success - corporation credentials retrieved", func() {
		corporationID := uint(1)
		corporation := &entity.Corporation{
			Name:   "Test Corporation",
			Status: enum.CorpStatusApproved,
		}

		addresses := []addressdto.AddressResponse{
			{ID: 1, ProvinceID: 1, CityID: 1, PostalCode: "12345"},
		}

		contactInfoEntities := []*entity.ContactInformation{
			{Value: "test@example.com", TypeID: 1},
		}

		contactType := &entity.ContactType{
			Name: "Email",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(corporation, nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return(contactInfoEntities, nil).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, uint(1)).Return(contactType, nil).Once()

		result, err := s.corporationService.GetCorporationCredentials(corporationID)

		s.NoError(err)
		s.Equal("Test Corporation", result.Name)
		s.Equal(enum.CorpStatusApproved.String(), result.Status)
		s.Len(result.Addresses, 1)
		s.Len(result.ContactInfo, 1)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		corporationID := uint(1)
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(nilCorporation, nil).Once()

		result, err := s.corporationService.GetCorporationCredentials(corporationID)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - address service error", func() {
		corporationID := uint(1)
		corporation := &entity.Corporation{
			Name:   "Test Corporation",
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, corporationID).Return(corporation, nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return([]addressdto.AddressResponse{}, errors.New("address service error")).Once()

		result, err := s.corporationService.GetCorporationCredentials(corporationID)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestAddContactInfo() {
	s.Run("success - contact info added", func() {
		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
			ContactInformation: []corporationdto.ContactInformation{
				{ContactTypeID: 1, ContactValue: "test@example.com"},
			},
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		staff := &entity.CorporationStaff{
			UserID:        2,
			CorporationID: 1,
		}

		contactType := &entity.ContactType{
			Name: "Email",
		}

		var nilContact *entity.ContactInformation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(staff, nil).Once()
		s.corporationRepository.On("FindContactInformationTypeByID", s.db, uint(1)).Return(contactType, nil).Once()
		s.corporationRepository.On("FindContactInformationTypeValue", s.db, uint(1), "test@example.com").Return(nilContact, nil).Once()
		s.corporationRepository.On("CreateContactInformation", s.db, mock.AnythingOfType("*entity.ContactInformation")).Return(nil).Once()

		err := s.corporationService.AddContactInfo(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - user not active", func() {
		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(errors.New("user not active")).Once()

		err := s.corporationService.AddContactInfo(request)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - applicant access denied", func() {
		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		var nilStaff *entity.CorporationStaff = nil

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(nilStaff, nil).Once()

		err := s.corporationService.AddContactInfo(request)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestDeleteContactInfo() {
	s.Run("success - contact info deleted", func() {
		request := corporationdto.DeleteContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
			ContactID:         3,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		staff := &entity.CorporationStaff{
			UserID:        2,
			CorporationID: 1,
		}

		contact := &entity.ContactInformation{
			CorporationID: 1,
			Value:         "test@example.com",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(staff, nil).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(3)).Return(contact, nil).Once()
		s.corporationRepository.On("DeleteContactInfo", s.db, contact).Return(nil).Once()

		err := s.corporationService.DeleteContactInfo(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - contact not found", func() {
		request := corporationdto.DeleteContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
			ContactID:         3,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		staff := &entity.CorporationStaff{
			UserID:        2,
			CorporationID: 1,
		}

		var nilContact *entity.ContactInformation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(staff, nil).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(3)).Return(nilContact, nil).Once()

		err := s.corporationService.DeleteContactInfo(request)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - delete operation failed", func() {
		request := corporationdto.DeleteContactInformationRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
			ContactID:         3,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		staff := &entity.CorporationStaff{
			UserID:        2,
			CorporationID: 1,
		}

		contact := &entity.ContactInformation{
			CorporationID: 1,
			Value:         "test@example.com",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(staff, nil).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(3)).Return(contact, nil).Once()
		s.corporationRepository.On("DeleteContactInfo", s.db, contact).Return(errors.New("delete failed")).Once()

		err := s.corporationService.DeleteContactInfo(request)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestRegister() {
	s.Run("success - corporation registered", func() {
		request := corporationdto.RegisterRequest{
			ApplicantID:        1,
			Name:               "Test Corporation",
			RegistrationNumber: "REG123",
			NationalID:         "NAT123",
		}

		role := &entity.Role{
			Name: enum.CorporationOwner.String(),
		}

		var nilCorporation *entity.Corporation = nil

		s.userService.On("IsUserActive", uint(1)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "Test Corporation", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "NAT123", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "REG123", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindRoleByName", s.db, enum.CorporationOwner.String()).Return(role, nil).Once()
		s.corporationRepository.On("CreateCorporation", mock.Anything, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("CreateCorporationStaff", mock.Anything, mock.Anything).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		result, err := s.corporationService.Register(request)

		s.NoError(err)
		s.NotEqual(corporationdto.CorporationCredentialResponse{}, result)
		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - user not active", func() {
		request := corporationdto.RegisterRequest{
			ApplicantID: 1,
		}

		s.userService.On("IsUserActive", uint(1)).Return(errors.New("user not active")).Once()

		result, err := s.corporationService.Register(request)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - duplicate corporation name", func() {
		request := corporationdto.RegisterRequest{
			ApplicantID: 1,
			Name:        "Existing Corporation",
		}

		existingCorporation := &entity.Corporation{
			Name: "Existing Corporation",
		}
		var nilCorporation *entity.Corporation = nil

		s.userService.On("IsUserActive", uint(1)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "Existing Corporation", mock.Anything).Return(existingCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "", mock.Anything).Return(nilCorporation, nil).Once()

		result, err := s.corporationService.Register(request)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - role not found", func() {
		request := corporationdto.RegisterRequest{
			ApplicantID:        1,
			Name:               "Test Corporation",
			RegistrationNumber: "REG123",
			NationalID:         "NAT123",
		}

		var nilRole *entity.Role = nil
		var nilCorporation *entity.Corporation = nil
		s.userService.On("IsUserActive", uint(1)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "Test Corporation", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "NAT123", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "REG123", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindRoleByName", s.db, enum.CorporationOwner.String()).Return(nilRole, nil).Once()

		result, err := s.corporationService.Register(request)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestUpdateRegistrationInfoProfile() {
	s.Run("success - profile updated", func() {
		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   2,
			Name:          &[]string{"Updated Corp"}[0],
			NationalID:    &[]string{"NEW123"}[0],
		}

		corporation := &entity.Corporation{
			Name:       "Old Corp",
			NationalID: "OLD123",
			Status:     enum.CorpStatusApproved,
		}

		var nilCorporation *entity.Corporation
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "Updated Corp", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "NEW123", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.corporationRepository.On("DeleteCorporationSignatories", mock.Anything, uint(1)).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.UpdateRegistrationInfoProfile(request)

		s.NoError(err)
		s.Equal("Updated Corp", corporation.Name)
		s.Equal("NEW123", corporation.NationalID)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		var nilCorporation *entity.Corporation
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(nilCorporation, nil).Once()

		err := s.corporationService.UpdateRegistrationInfoProfile(request)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - user not active", func() {
		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(errors.New("user not active")).Once()

		err := s.corporationService.UpdateRegistrationInfoProfile(request)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestUpdateRegister() {
	s.Run("success - registration updated", func() {
		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   2,
			Name:          &[]string{"Updated Corp"}[0],
		}

		corporation := &entity.Corporation{
			Name:   "Old Corp",
			Status: enum.CorpStatusAwaitingApproval,
		}

		var nilCorporation *entity.Corporation
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "Updated Corp", mock.Anything).Return(nilCorporation, nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.corporationRepository.On("DeleteCorporationSignatories", mock.Anything, uint(1)).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.UpdateRegister(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - corporation status not awaiting approval", func() {
		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()

		err := s.corporationService.UpdateRegister(request)

		s.Error(err)
		s.IsType(exception.ForbiddenError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestAddCertificateFilesFromProfile() {
	s.Run("success - certificates added from profile", func() {
		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(nil).Once()

		err := s.corporationService.AddCertificateFilesFromProfile(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		var nilCorporation *entity.Corporation
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(nilCorporation, nil).Once()

		err := s.corporationService.AddCertificateFilesFromProfile(request)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestAddCertificateFiles() {
	s.Run("success - certificates added", func() {
		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(nil).Once()

		err := s.corporationService.AddCertificateFiles(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - corporation status not awaiting approval", func() {
		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.AddCertificateFiles(request)

		s.Error(err)
		s.IsType(exception.ForbiddenError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationDetails() {
	s.Run("success - corporation details retrieved", func() {
		request := corporationdto.CorporationDetailsRequest{
			CorporationID: 1,
			UserID:        2,
		}

		corporation := &entity.Corporation{
			Name:   "Test Corp",
			Status: enum.CorpStatusApproved,
		}

		addresses := []addressdto.AddressResponse{
			{ID: 1, ProvinceID: 1, CityID: 1},
		}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return([]*entity.ContactInformation{}, nil).Once()
		s.corporationRepository.On("FindCorporationSignatories", s.db, mock.Anything).Return([]*entity.Signatory{}, nil).Once()

		result, err := s.corporationService.GetCorporationDetails(request)

		s.NoError(err)
		s.Equal("Test Corp", result.Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - access denied", func() {
		request := corporationdto.CorporationDetailsRequest{
			CorporationID: 1,
			UserID:        2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		var nilStaff *entity.CorporationStaff
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(nilStaff, nil).Once()

		result, err := s.corporationService.GetCorporationDetails(request)

		s.Error(err)
		s.Equal(corporationdto.CorporationPrivateInfoResponse{}, result)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationPublicDetails() {
	s.Run("success - public details retrieved", func() {
		request := corporationdto.CorporationDetailsRequest{
			CorporationID: 1,
			UserID:        2,
		}

		corporation := &entity.Corporation{
			Name:   "Test Corp",
			Status: enum.CorpStatusApproved,
		}

		addresses := []addressdto.AddressResponse{}

		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, mock.Anything).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return([]*entity.ContactInformation{}, nil).Once()

		result, err := s.corporationService.GetCorporationPublicDetails(request)

		s.NoError(err)
		s.Equal("Test Corp", result.Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - access denied", func() {
		request := corporationdto.CorporationDetailsRequest{
			CorporationID: 1,
			UserID:        2,
		}

		var nilStaff *entity.CorporationStaff
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(nilStaff, nil).Once()

		result, err := s.corporationService.GetCorporationPublicDetails(request)

		s.Error(err)
		s.Equal(corporationdto.CorporationCredentialResponse{}, result)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetContactTypes() {
	s.Run("success - contact types retrieved", func() {
		contactTypes := []*entity.ContactType{
			{Name: "Email"},
			{Name: "Phone"},
		}

		s.corporationRepository.On("FindContactTypes", s.db).Return(contactTypes, nil).Once()

		result, err := s.corporationService.GetContactTypes()

		s.NoError(err)
		s.Len(result, 2)
		s.Equal("Email", result[0].Name)
		s.Equal("Phone", result[1].Name)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		var nilContactTypes []*entity.ContactType
		s.corporationRepository.On("FindContactTypes", s.db).Return(nilContactTypes, errors.New("database error")).Once()

		result, err := s.corporationService.GetContactTypes()

		s.Error(err)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestAddAddress() {
	s.Run("success - address added", func() {
		request := corporationdto.AddCorporationAddressRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       2,
			Addresses:         []addressdto.CreateAddressRequest{{PostalCode: "12345"}},
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(&entity.CorporationStaff{}, nil).Once()
		s.addressService.On("CreateAddress", mock.Anything).Return(addressdto.AddressResponse{}, nil).Once()

		err := s.corporationService.AddAddress(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - corporation status mismatch", func() {
		request := corporationdto.AddCorporationAddressRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusAwaitingApproval,
			ApplicantID:       2,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.AddAddress(request)

		s.Error(err)
		s.IsType(exception.ForbiddenError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestDeleteAddress() {
	s.Run("success - address deleted", func() {
		request := corporationdto.DeleteAddressRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			UserID:            2,
			AddressID:         3,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(&entity.CorporationStaff{}, nil).Once()
		s.addressService.On("DeleteAddress", uint(3)).Return(nil).Once()

		err := s.corporationService.DeleteAddress(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - address service error", func() {
		request := corporationdto.DeleteAddressRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			UserID:            2,
			AddressID:         3,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusApproved}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.userService.On("IsUserActive", uint(2)).Return(nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(&entity.CorporationStaff{}, nil).Once()
		s.addressService.On("DeleteAddress", uint(3)).Return(errors.New("delete failed")).Once()

		err := s.corporationService.DeleteAddress(request)

		s.Error(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestChangeLogo() {
	s.Run("success - logo changed", func() {
		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(2), uint(1)).Return(&entity.CorporationStaff{}, nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.ChangeLogo(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - corporation not approved", func() {
		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   2,
		}

		corporation := &entity.Corporation{Status: enum.CorpStatusAwaitingApproval}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.ChangeLogo(request)

		s.Error(err)
		s.IsType(exception.ForbiddenError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetUserCorporations() {
	s.Run("success - user corporations retrieved", func() {
		userID := uint(1)
		corporations := []*entity.Corporation{
			{Name: "Corp 1", Status: enum.CorpStatusApproved},
			{Name: "Corp 2", Status: enum.CorpStatusApproved},
		}

		addresses := []addressdto.AddressResponse{}

		s.corporationRepository.On("FindUserActiveCorporations", s.db, userID).Return(corporations, nil).Once()
		for _, corp := range corporations {
			s.corporationRepository.On("FindCorporationByID", s.db, corp.ID).Return(corp, nil).Once()
			s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
			s.corporationRepository.On("FindContactInformation", s.db, corp.ID).Return([]*entity.ContactInformation{}, nil).Once()
		}

		result, err := s.corporationService.GetUserCorporations(userID)

		s.NoError(err)
		s.Len(result, 2)
		s.Equal("Corp 1", result[0].Name)
		s.Equal("Corp 2", result[1].Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		userID := uint(1)

		var nilCorporations []*entity.Corporation
		s.corporationRepository.On("FindUserActiveCorporations", s.db, userID).Return(nilCorporations, errors.New("database error")).Once()

		result, err := s.corporationService.GetUserCorporations(userID)

		s.Error(err)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetAvailableCorporations() {
	s.Run("success - available corporations retrieved", func() {
		corporations := []*entity.Corporation{
			{Name: "Corp 1", Status: enum.CorpStatusApproved},
		}

		addresses := []addressdto.AddressResponse{}

		s.corporationRepository.On("FindCorporationsByStatus", s.db, []enum.CorporationStatus{enum.CorpStatusApproved}, (*postgres.QueryOptions)(nil)).Return(corporations, nil).Once()
		s.corporationRepository.On("FindCorporationByID", s.db, uint(0)).Return(corporations[0], nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, uint(0)).Return([]*entity.ContactInformation{}, nil).Once()

		result, err := s.corporationService.GetAvailableCorporations()

		s.NoError(err)
		s.Len(result, 1)
		s.Equal("Corp 1", result[0].Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		var nilCorporations []*entity.Corporation
		s.corporationRepository.On("FindCorporationsByStatus", s.db, mock.Anything, mock.Anything).Return(nilCorporations, errors.New("database error")).Once()

		result, err := s.corporationService.GetAvailableCorporations()

		s.Error(err)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetReviewActions() {
	s.Run("success - review actions retrieved", func() {
		result := s.corporationService.GetReviewActions()

		s.NotNil(result)
		s.Greater(len(result), 0)

		for _, action := range result {
			s.Greater(action.ID, uint(0))
			s.NotEmpty(action.Name)
		}
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationReviewsByAdmin() {
	s.Run("success - reviews retrieved", func() {
		corporationID := uint(1)
		reviews := []*entity.CorporationReview{
			{
				ReviewerID: 1,
				Action:     enum.ReviewActionApproved,
				Reason:     stringPtr("Approved"),
				Notes:      stringPtr("All good"),
			},
		}

		userCredential := userdto.CredentialResponse{
			Phone: "123456789",
		}

		s.corporationRepository.On("FindCorporationReviews", s.db, corporationID, (*postgres.QueryOptions)(nil)).Return(reviews, nil).Once()
		s.userService.On("GetUserCredential", uint(1)).Return(userCredential, nil).Once()

		result, err := s.corporationService.GetCorporationReviewsByAdmin(corporationID)

		s.NoError(err)
		s.Len(result, 1)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		corporationID := uint(1)

		var nilReviews []*entity.CorporationReview
		s.corporationRepository.On("FindCorporationReviews", s.db, corporationID, (*postgres.QueryOptions)(nil)).Return(nilReviews, errors.New("database error")).Once()

		result, err := s.corporationService.GetCorporationReviewsByAdmin(corporationID)

		s.Error(err)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestApproveCorporationRegistration() {
	s.Run("success - corporation approved", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
			Reason:        stringPtr("Approved"),
			Notes:         stringPtr("All documents verified"),
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.corporationRepository.On("CreateReview", mock.Anything, mock.AnythingOfType("*entity.CorporationReview")).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.ApproveCorporationRegistration(request)

		s.NoError(err)
		s.Equal(enum.CorpStatusApproved, corporation.Status)
		s.corporationRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - corporation already approved", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.ApproveCorporationRegistration(request)

		s.Error(err)
		s.IsType(exception.ConflictErrors{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - corporation already rejected", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusRejected,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.ApproveCorporationRegistration(request)

		s.Error(err)
		s.IsType(exception.ConflictErrors{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestRejectCorporationRegistration() {
	s.Run("success - corporation rejected", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
			ActionID:      uint(enum.ReviewActionRejected),
			Reason:        stringPtr("Incomplete documents"),
			Notes:         stringPtr("Missing tax certificate"),
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.corporationRepository.On("CreateReview", mock.Anything, mock.AnythingOfType("*entity.CorporationReview")).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.RejectCorporationRegistration(request)

		s.NoError(err)
		s.Equal(enum.CorpStatusRejected, corporation.Status)
		s.corporationRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("success - corporation suspended", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
			ActionID:      uint(enum.ReviewActionSuspended),
			Reason:        stringPtr("Suspicious activity"),
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()
		s.corporationRepository.On("CreateReview", mock.Anything, mock.AnythingOfType("*entity.CorporationReview")).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", mock.Anything, corporation).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.RejectCorporationRegistration(request)

		s.NoError(err)
		s.Equal(enum.CorpStatusSuspend, corporation.Status)
		s.corporationRepository.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - approval action not allowed", func() {
		request := corporationdto.HandleCorporationActionRequest{
			CorporationID: 1,
			ReviewerID:    2,
			ActionID:      uint(enum.ReviewActionApproved),
		}

		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, nil).Once()

		err := s.corporationService.RejectCorporationRegistration(request)

		s.Error(err)
		s.IsType(exception.ForbiddenError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetStaffStatuses() {
	s.Run("success - staff statuses retrieved", func() {
		result := s.corporationService.GetStaffStatuses()

		s.NotNil(result)
		s.Greater(len(result), 0)

		// Verify each status has ID and Name
		for _, status := range result {
			s.Greater(status.ID, uint(0))
			s.NotEmpty(status.Name)
		}
	})
}

func (s *CorporationServiceTestSuite) TestAddStaff() {
	s.Run("success - new staff added", func() {
		request := corporationdto.AddStaffRequest{
			CorporationID: 1,
			StaffPhone:    "123456789",
			RoleIDs:       []uint{1, 2},
		}

		user := &entity.User{}

		var nilStaff *entity.CorporationStaff

		s.userService.On("FindActiveUserByPhone", "123456789").Return(user, nil).Once()
		s.corporationRepository.On("FindStaffByUserIDAndStatus", s.db, mock.Anything, mock.Anything).Return(nilStaff, nil).Once()
		s.corporationRepository.On("CreateStaff", s.db, mock.Anything).Return(nil).Once()
		s.rbacService.On("UpdateStaffRoles", mock.Anything).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.AddStaff(request)

		s.NoError(err)
		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("success - existing staff roles updated", func() {
		request := corporationdto.AddStaffRequest{
			CorporationID: 1,
			StaffPhone:    "123456789",
			RoleIDs:       []uint{1, 2},
		}

		user := &entity.User{}

		staff := &entity.CorporationStaff{
			CorporationID: 1,
			Status:        enum.StaffStatusActive,
		}

		s.userService.On("FindActiveUserByPhone", "123456789").Return(user, nil).Once()
		s.corporationRepository.On("FindStaffByUserIDAndStatus", s.db, mock.Anything, mock.Anything).Return(staff, nil).Once()
		s.rbacService.On("UpdateStaffRoles", mock.Anything).Return(nil).Once()

		err := s.corporationService.AddStaff(request)

		s.NoError(err)
		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})

	s.Run("error - user not found", func() {
		request := corporationdto.AddStaffRequest{
			CorporationID: 1,
			StaffPhone:    "123456789",
			RoleIDs:       []uint{1, 2},
		}

		var nilUser *entity.User
		s.userService.On("FindActiveUserByPhone", "123456789").Return(nilUser, errors.New("user not found")).Once()

		err := s.corporationService.AddStaff(request)

		s.Error(err)
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestEditStaff() {
	s.Run("success - staff edited", func() {
		request := corporationdto.EditStaffRequest{
			CorporationID: 1,
			StaffID:       2,
			RoleIDs:       []uint{1, 3},
			Status:        &[]uint{uint(enum.StaffStatusActive)}[0],
		}

		staff := &entity.CorporationStaff{
			UserID:        3,
			CorporationID: 1,
			Status:        enum.StaffStatusActive,
		}

		s.corporationRepository.On("FindCorporationStaffByID", s.db, uint(1), uint(2)).Return(staff, nil).Once()
		s.rbacService.On("UpdateStaffRoles", mock.AnythingOfType("rbacdto.UpdateStaffRolesRequest")).Return(nil).Once()
		s.corporationRepository.On("UpdateStaff", s.db, staff).Return(nil).Once()
		s.db.On("WithTransaction", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(database.Database) error)
			fn(s.db)
		})

		err := s.corporationService.EditStaff(request)

		s.NoError(err)
		s.corporationRepository.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("error - staff not active", func() {
		request := corporationdto.EditStaffRequest{
			CorporationID: 1,
			StaffID:       2,
			RoleIDs:       []uint{1, 3},
		}

		staff := &entity.CorporationStaff{
			UserID:        3,
			CorporationID: 1,
			Status:        enum.StaffStatusInactive,
		}

		s.corporationRepository.On("FindCorporationStaffByID", s.db, uint(1), uint(2)).Return(staff, nil).Once()

		err := s.corporationService.EditStaff(request)

		s.Error(err)
		s.IsType(exception.ConflictErrors{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})

	s.Run("error - staff not found", func() {
		request := corporationdto.EditStaffRequest{
			CorporationID: 1,
			StaffID:       2,
			RoleIDs:       []uint{1, 3},
		}

		var nilStaff *entity.CorporationStaff
		s.corporationRepository.On("FindCorporationStaffByID", s.db, uint(1), uint(2)).Return(nilStaff, nil).Once()

		err := s.corporationService.EditStaff(request)

		s.Error(err)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetStaffList() {
	s.Run("success - staff list retrieved", func() {
		request := corporationdto.GetStaffList{
			CorporationID: 1,
			Status:        uint(enum.StaffStatusAll),
			Query:         "",
			SortBy:        1,
			Asc:           true,
			Limit:         10,
			Offset:        0,
		}

		staffs := []*entity.CorporationStaff{
			{
				UserID:        2,
				CorporationID: 1,
				Status:        enum.StaffStatusActive,
			},
		}

		userCredential := userdto.CredentialResponse{
			Phone: "123456789",
		}

		roles := []rbacdto.RoleResponse{
			{Name: "Admin"},
		}

		s.corporationRepository.On("FindCorporationStaffs", s.db, uint(1), mock.Anything, mock.Anything).Return(staffs, nil).Once()
		s.corporationRepository.On("CountCorporationStaffs", s.db, uint(1), mock.Anything).Return(int64(1), nil).Once()
		s.userService.On("GetUserCredential", uint(2)).Return(userCredential, nil).Once()
		s.rbacService.On("GetStaffRoles", staffs[0]).Return(roles, nil).Once()

		result, count, err := s.corporationService.GetStaffList(request)

		s.NoError(err)
		s.Equal(int64(1), count)
		s.Len(result, 1)
		s.Equal("123456789", result[0].Staff.Phone)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		request := corporationdto.GetStaffList{
			CorporationID: 1,
			Status:        uint(enum.StaffStatusAll),
		}

		var nilStaffs []*entity.CorporationStaff
		s.corporationRepository.On("FindCorporationStaffs", s.db, uint(1), mock.Anything, mock.Anything).Return(nilStaffs, errors.New("database error")).Once()

		result, count, err := s.corporationService.GetStaffList(request)

		s.Error(err)
		s.Equal(int64(0), count)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetStaff() {
	s.Run("success - staff retrieved", func() {
		corporationID := uint(1)
		staffID := uint(2)

		staff := &entity.CorporationStaff{
			UserID:        3,
			CorporationID: 1,
			Status:        enum.StaffStatusActive,
		}

		userCredential := userdto.CredentialResponse{
			Phone: "123456789",
		}

		roles := []rbacdto.RoleResponse{
			{Name: "Admin"},
		}

		s.corporationRepository.On("FindCorporationStaffByID", s.db, corporationID, staffID).Return(staff, nil).Once()
		s.userService.On("GetUserCredential", uint(3)).Return(userCredential, nil).Once()
		s.rbacService.On("GetStaffRoles", staff).Return(roles, nil).Once()

		result, err := s.corporationService.GetStaff(corporationID, staffID)

		s.NoError(err)
		s.Equal("123456789", result.Staff.Phone)
		s.Len(result.Roles, 1)
		s.Equal("Admin", result.Roles[0].Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.rbacService.AssertExpectations(s.T())
	})

	s.Run("error - staff not found", func() {
		corporationID := uint(1)
		staffID := uint(2)

		var nilStaff *entity.CorporationStaff
		s.corporationRepository.On("FindCorporationStaffByID", s.db, corporationID, staffID).Return(nilStaff, nil).Once()

		result, err := s.corporationService.GetStaff(corporationID, staffID)

		s.Error(err)
		s.Equal(corporationdto.StaffDetailsResponse{}, result)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationsByAdmin() {
	s.Run("success - corporations retrieved by admin", func() {
		request := corporationdto.GetCorporationsByAdminRequest{
			Status: uint(enum.CorpStatusAll),
			Query:  "",
			SortBy: 1,
			Asc:    true,
			Limit:  10,
			Offset: 0,
		}

		corporations := []*entity.Corporation{
			{Name: "Corp 1", Status: enum.CorpStatusApproved},
		}

		addresses := []addressdto.AddressResponse{}

		s.corporationRepository.On("FindCorporationsByStatus", s.db, mock.Anything, mock.Anything).Return(corporations, nil).Once()
		s.corporationRepository.On("CountCorporationsByStatus", s.db, mock.Anything).Return(int64(1), nil).Once()
		s.corporationRepository.On("FindCorporationByID", s.db, uint(0)).Return(corporations[0], nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, uint(0)).Return([]*entity.ContactInformation{}, nil).Once()

		result, count, err := s.corporationService.GetCorporationsByAdmin(request)

		s.NoError(err)
		s.Equal(int64(1), count)
		s.Len(result, 1)
		s.Equal("Corp 1", result[0].Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - repository error", func() {
		request := corporationdto.GetCorporationsByAdminRequest{
			Status: uint(enum.CorpStatusAll),
		}

		var nilCorporations []*entity.Corporation
		s.corporationRepository.On("FindCorporationsByStatus", s.db, mock.Anything, mock.Anything).Return(nilCorporations, errors.New("database error")).Once()

		result, count, err := s.corporationService.GetCorporationsByAdmin(request)

		s.Error(err)
		s.Equal(int64(0), count)
		s.Nil(result)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationByAdmin() {
	s.Run("success - corporation retrieved by admin", func() {
		corporationID := uint(1)

		corporation := &entity.Corporation{
			Name:   "Test Corp",
			Status: enum.CorpStatusApproved,
		}

		addresses := []addressdto.AddressResponse{}

		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(corporation, nil).Once()
		s.addressService.On("GetAddresses", mock.AnythingOfType("addressdto.GetOwnerAddressesRequest")).Return(addresses, nil).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return([]*entity.ContactInformation{}, nil).Once()
		s.corporationRepository.On("FindCorporationSignatories", s.db, mock.Anything).Return([]*entity.Signatory{}, nil).Once()

		result, err := s.corporationService.GetCorporationByAdmin(corporationID)

		s.NoError(err)
		s.Equal("Test Corp", result.Name)
		s.corporationRepository.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("error - corporation not found", func() {
		corporationID := uint(1)

		var nilCorporation *entity.Corporation
		s.corporationRepository.On("FindCorporationByID", s.db, mock.Anything).Return(nilCorporation, nil).Once()

		result, err := s.corporationService.GetCorporationByAdmin(corporationID)

		s.Error(err)
		s.Equal(corporationdto.CorporationPrivateInfoResponse{}, result)
		s.IsType(exception.NotFoundError{}, err)
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func TestCorporationService(t *testing.T) {
	suite.Run(t, new(CorporationServiceTestSuite))
}
