package serviceimpl

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
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

func (s *CorporationServiceTestSuite) TestRegister() {
	s.Run("success - Corporation registered", func() {
		var nilCorporation *entity.Corporation = nil
		var nilSignatory *entity.Signatory = nil
		signatory := &entity.Signatory{}

		s.userService.On("IsUserActive", mock.Anything).Return(true).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "testName", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "testNationalID", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "testRegistrationNumber", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("CreateCorporation", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("CreateCorporationStaff", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(signatory, true).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(nilSignatory, false).Once()
		s.corporationRepository.On("CreateSignatory", s.db, mock.Anything).Return(nil).Once()

		request := corporationdto.RegisterRequest{
			ApplicantID:        1,
			Name:               "testName",
			NationalID:         "testNationalID",
			RegistrationNumber: "testRegistrationNumber",
			IBAN:               "testIBAN",
			Signatories:        []corporationdto.Signatory{{}, {}},
		}

		response := s.corporationService.Register(request)

		s.Equal(response.Name, request.Name)

		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		s.userService.On("IsUserActive", mock.Anything).Return(false).Once()

		s.Panics(func() {
			s.corporationService.Register(corporationdto.RegisterRequest{})
		})

		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Corporation credentials already exist", func() {
		corporation := &entity.Corporation{
			Name:               "testName",
			NationalID:         "testNationalID",
			IBAN:               "testIBAN",
			RegistrationNumber: "testRegistrationNumber",
		}

		s.userService.On("IsUserActive", mock.Anything).Return(true).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "testName", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "testNationalID", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "testRegistrationNumber", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(corporation, true).Once()

		request := corporationdto.RegisterRequest{
			Name:               "testName",
			NationalID:         "testNationalID",
			RegistrationNumber: "testRegistrationNumber",
			IBAN:               "testIBAN",
		}

		s.Panics(func() {
			s.corporationService.Register(request)
		})

		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Create corporation failed", func() {
		var nilCorporation *entity.Corporation = nil

		s.userService.On("IsUserActive", mock.Anything).Return(true).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()

		s.corporationRepository.On("CreateCorporation", s.db, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.RegisterRequest{
			IBAN: "testIBAN",
		}

		s.Panics(func() {
			s.corporationService.Register(request)
		})

		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Create corporation staff failed", func() {
		var nilCorporation *entity.Corporation = nil

		s.userService.On("IsUserActive", mock.Anything).Return(true).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("CreateCorporation", s.db, mock.Anything).Return(nil).Once()

		s.corporationRepository.On("CreateCorporationStaff", s.db, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.RegisterRequest{
			IBAN: "testIBAN",
		}

		s.Panics(func() {
			s.corporationService.Register(request)
		})

		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Create signatory failed", func() {
		var nilCorporation *entity.Corporation = nil
		var nilSignatory *entity.Signatory = nil
		signatory := &entity.Signatory{}

		s.userService.On("IsUserActive", mock.Anything).Return(true).Once()
		s.corporationRepository.On("FindCorporationByName", s.db, "testName", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "testNationalID", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "testRegistrationNumber", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("CreateCorporation", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("CreateCorporationStaff", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(signatory, true).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(nilSignatory, false).Once()
		s.corporationRepository.On("CreateSignatory", s.db, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.RegisterRequest{
			ApplicantID:        1,
			Name:               "testName",
			NationalID:         "testNationalID",
			RegistrationNumber: "testRegistrationNumber",
			IBAN:               "testIBAN",
			Signatories:        []corporationdto.Signatory{{}, {}},
		}

		s.Panics(func() {
			s.corporationService.Register(request)
		})

		s.userService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestReplaceSignatories() {
	s.Run("success - Signatories replaced", func() {
		corporation := &entity.Corporation{
			Signatories: []entity.Signatory{
				{
					NationalCardNumber: "1234567890",
					Position:           "existingPosition",
				},
			},
		}
		var nilSignatory *entity.Signatory = nil

		s.corporationRepository.On("DeleteCorporationSignatories", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(nilSignatory, false).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(&corporation.Signatories[0], true).Once()
		s.corporationRepository.On("CreateSignatory", s.db, mock.Anything).Return(nil).Once()

		s.corporationService.replaceSignatories(uint(1), []corporationdto.Signatory{{
			NationalCardNumber: "1234567890",
			Position:           "existingPosition",
		}, {
			NationalCardNumber: "1234567891",
			Position:           "newPosition",
		}})

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Delete corporation signatories failed", func() {
		s.corporationRepository.On("DeleteCorporationSignatories", s.db, mock.Anything).Return(errors.New("error")).Once()

		s.Panics(func() {
			s.corporationService.replaceSignatories(uint(1), []corporationdto.Signatory{{}, {}})
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Create signatory failed", func() {
		var nilSignatory *entity.Signatory = nil

		s.corporationRepository.On("DeleteCorporationSignatories", s.db, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("FindCorporationSignatoryByNationalID", s.db, mock.Anything, mock.Anything, mock.Anything).Return(nilSignatory, false).Once()
		s.corporationRepository.On("CreateSignatory", s.db, mock.Anything).Return(errors.New("error")).Once()

		s.Panics(func() {
			s.corporationService.replaceSignatories(uint(1), []corporationdto.Signatory{{}})
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestUpdateRegister() {
	s.Run("success - Corporation updated", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(nil).Once()
		s.corporationRepository.On("DeleteCorporationSignatories", s.db, mock.Anything).Return(nil).Once()

		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   1,
		}
		s.corporationService.UpdateRegister(request)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   1,
		}
		s.Panics(func() {
			s.corporationService.UpdateRegister(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Update corporation failed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(errors.New("error")).Once()

		request := corporationdto.UpdateRegisterRequest{
			CorporationID: 1,
			ApplicantID:   1,
		}

		s.Panics(func() {
			s.corporationService.UpdateRegister(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestCheckCorporationConflicts() {
	s.Run("success - Corporation conflicts checked", func() {
		corporation := &entity.Corporation{
			Name:               "testName",
			NationalID:         "testNationalID",
			RegistrationNumber: "testRegistrationNumber",
			IBAN:               "testIBAN",
		}
		var nilCorporation *entity.Corporation = nil

		s.corporationRepository.On("FindCorporationByName", s.db, "testName", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "testNationalID", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "testRegistrationNumber", mock.Anything).Return(nilCorporation, false).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(nilCorporation, false).Once()

		s.corporationService.checkCorporationConflicts(corporation, &corporation.Name, &corporation.NationalID, &corporation.RegistrationNumber, &corporation.IBAN)

		s.corporationRepository.AssertExpectations(s.T())
	})
	s.Run("error - Corporation conflicts checked", func() {
		corporation := &entity.Corporation{
			Name:               "testName",
			NationalID:         "testNationalID",
			RegistrationNumber: "testRegistrationNumber",
			IBAN:               "testIBAN",
		}

		s.corporationRepository.On("FindCorporationByName", s.db, "testName", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByNationalID", s.db, "testNationalID", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByRegistrationNumber", s.db, "testRegistrationNumber", mock.Anything).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationByIBAN", mock.Anything, mock.Anything, mock.Anything).Return(corporation, true).Once()

		s.Panics(func() {
			s.corporationService.checkCorporationConflicts(corporation, &corporation.Name, &corporation.NationalID, &corporation.RegistrationNumber, &corporation.IBAN)
		})

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestAddCertificateFiles() {
	s.Run("success - Certificate files added", func() {
		corporationStaff := &entity.CorporationStaff{}
		corporation := &entity.Corporation{
			Status:                 enum.CorpStatusAwaitingApproval,
			VATTaxpayerCertificate: "testVATTaxpayerCertificate",
			OfficialNewspaperAD:    "testOfficialNewspaperAD",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.VATTaxpayerCertificate, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.VATTaxpayerCertificate, mock.Anything).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.OfficialNewspaperAD, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.OfficialNewspaperAD, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(nil).Once()

		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   1,
			VATTaxpayerCertificate: &multipart.FileHeader{
				Filename: "testVATTaxpayerCertificate",
			},
			OfficialNewspaperAD: &multipart.FileHeader{
				Filename: "testOfficialNewspaperAD",
			},
		}
		s.corporationService.AddCertificateFiles(request)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusAwaitingApproval,
		}
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   1,
		}
		s.Panics(func() {
			s.corporationService.AddCertificateFiles(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Delete VAT Taxpayer Certificate failed", func() {
		corporation := &entity.Corporation{
			Status:                 enum.CorpStatusAwaitingApproval,
			VATTaxpayerCertificate: "testVATTaxpayerCertificate",
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.VATTaxpayerCertificate, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.VATTaxpayerCertificate, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   1,
			VATTaxpayerCertificate: &multipart.FileHeader{
				Filename: "testVATTaxpayerCertificate",
			},
		}
		s.Panics(func() {
			s.corporationService.AddCertificateFiles(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
	})
	s.Run("error - Delete Official Newspaper AD failed", func() {
		corporation := &entity.Corporation{
			Status:              enum.CorpStatusAwaitingApproval,
			OfficialNewspaperAD: "testOfficialNewspaperAD",
		}
		corporationStaff := &entity.CorporationStaff{}
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.OfficialNewspaperAD, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.OfficialNewspaperAD, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   1,
			OfficialNewspaperAD: &multipart.FileHeader{
				Filename: "testOfficialNewspaperAD",
			},
		}
		s.Panics(func() {
			s.corporationService.AddCertificateFiles(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Update corporation failed", func() {
		corporationStaff := &entity.CorporationStaff{}
		corporation := &entity.Corporation{
			Status:                 enum.CorpStatusAwaitingApproval,
			VATTaxpayerCertificate: "testVATTaxpayerCertificate",
			OfficialNewspaperAD:    "testOfficialNewspaperAD",
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.VATTaxpayerCertificate, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.VATTaxpayerCertificate, mock.Anything).Return(nil).Once()
		s.s3Storage.On("UploadObject", enum.OfficialNewspaperAD, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.OfficialNewspaperAD, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, corporation).Return(errors.New("error")).Once()

		request := corporationdto.AddCertificatesRequest{
			CorporationID: 1,
			ApplicantID:   1,
			VATTaxpayerCertificate: &multipart.FileHeader{
				Filename: "testVATTaxpayerCertificate",
			},
			OfficialNewspaperAD: &multipart.FileHeader{
				Filename: "testOfficialNewspaperAD",
			},
		}
		s.Panics(func() {
			s.corporationService.AddCertificateFiles(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())

	})
}

func TestCorporationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CorporationServiceTestSuite))
}
