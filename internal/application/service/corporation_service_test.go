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

func (s *CorporationServiceTestSuite) TestAddContactInfo() {
	s.Run("success - Contact info added", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}
		var nilContactInformation *entity.ContactInformation = nil
		var nilContactType *entity.ContactType = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()

		s.corporationRepository.On("FindContactInformationTypeByID", s.db, uint(1)).Return(nilContactType, true).Once()
		s.corporationRepository.On("FindContactInformationTypeValue", s.db, uint(1), "testContactValue").Return(nilContactInformation, false).Once()

		s.corporationRepository.On("FindContactInformationTypeByID", s.db, uint(2)).Return(nilContactType, false).Once()
		// s.corporationRepository.On("FindContactInformationTypeValue", s.db, uint(2), "testContactValue2").Return(nilContactInformation, false).Once()

		s.corporationRepository.On("FindContactInformationTypeByID", s.db, uint(3)).Return(nilContactType, true).Once()
		s.corporationRepository.On("FindContactInformationTypeValue", s.db, uint(3), "testContactValue3").Return(nilContactInformation, true).Once()

		s.corporationRepository.On("CreateContactInformation", s.db, mock.Anything).Return(nil).Once()

		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			ApplicantID:       1,
			CorporationStatus: enum.CorpStatusApproved,
			ContactInformation: []corporationdto.ContactInformation{
				{
					ContactTypeID: 1,
					ContactValue:  "testContactValue",
				},
				{
					ContactTypeID: 2,
				},
				{
					ContactTypeID: 3,
					ContactValue:  "testContactValue3",
				},
			},
		}
		s.corporationService.AddContactInfo(request)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			ApplicantID:       1,
			CorporationStatus: enum.CorpStatusApproved,
		}
		s.Panics(func() {
			s.corporationService.AddContactInfo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Create contact information failed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}
		var nilContactInformation *entity.ContactInformation = nil
		var nilContactType *entity.ContactType = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("FindContactInformationTypeByID", s.db, uint(1)).Return(nilContactType, true).Once()
		s.corporationRepository.On("FindContactInformationTypeValue", s.db, uint(1), "testContactValue").Return(nilContactInformation, false).Once()
		s.corporationRepository.On("CreateContactInformation", s.db, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.AddContactInformationRequest{
			CorporationID:     1,
			ApplicantID:       1,
			CorporationStatus: enum.CorpStatusApproved,
			ContactInformation: []corporationdto.ContactInformation{
				{
					ContactTypeID: 1,
					ContactValue:  "testContactValue",
				},
			},
		}
		s.Panics(func() {
			s.corporationService.AddContactInfo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestDeleteContactInfo() {
	s.Run("success - Contact info deleted", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}
		contactInformation := &entity.ContactInformation{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(1)).Return(contactInformation, true).Once()
		s.corporationRepository.On("DeleteContactInfo", s.db, mock.Anything).Return(nil).Once()

		request := corporationdto.DeleteContactInformationRequest{
			ApplicantID:       1,
			ContactID:         1,
			CorporationStatus: enum.CorpStatusApproved,
			CorporationID:     1,
		}
		s.corporationService.DeleteContactInfo(request)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.DeleteContactInformationRequest{
			ApplicantID:       1,
			ContactID:         1,
			CorporationStatus: enum.CorpStatusApproved,
			CorporationID:     1,
		}
		s.Panics(func() {
			s.corporationService.DeleteContactInfo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Contact information not found", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}
		var nilContactInformation *entity.ContactInformation = nil

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(1)).Return(nilContactInformation, false).Once()

		request := corporationdto.DeleteContactInformationRequest{
			ApplicantID:       1,
			ContactID:         1,
			CorporationStatus: enum.CorpStatusApproved,
			CorporationID:     1,
		}
		s.Panics(func() {
			s.corporationService.DeleteContactInfo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())

	})
	s.Run("error - Delete contact information failed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}
		contactInformation := &entity.ContactInformation{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.corporationRepository.On("FindContactInformationByID", s.db, uint(1)).Return(contactInformation, true).Once()
		s.corporationRepository.On("DeleteContactInfo", s.db, mock.Anything).Return(errors.New("error")).Once()

		request := corporationdto.DeleteContactInformationRequest{
			ApplicantID:       1,
			ContactID:         1,
			CorporationStatus: enum.CorpStatusApproved,
			CorporationID:     1,
		}
		s.Panics(func() {
			s.corporationService.DeleteContactInfo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationDetails() {
	s.Run("success - Corporation details fetched", func() {
		corporation := &entity.Corporation{
			Name:                   "testCorporation",
			RegistrationNumber:     "testRegistrationNumber",
			NationalID:             "testNationalID",
			IBAN:                   "testIBAN",
			Status:                 enum.CorpStatusApproved,
			VATTaxpayerCertificate: "testVATTaxpayerCertificate",
			OfficialNewspaperAD:    "testOfficialNewspaperAD",
			Logo:                   "testLogo",
		}
		contactInformation := []*entity.ContactInformation{
			{
				TypeID: 1,
				Value:  "testContactValue",
			},
			{
				TypeID: 2,
				Value:  "testContactValue2",
			},
		}
		contactType := &entity.ContactType{
			Name: "testContactType",
		}
		signatories := []*entity.Signatory{
			{
				Name: "testSignatory",
			},
		}
		var nilContactType *entity.ContactType = nil
		corporationStaff := &entity.CorporationStaff{}

		s.userService.On("DoesUserExist", mock.Anything).Return(nil).Once()
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, mock.Anything, uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("GetPresignedURL", enum.VATTaxpayerCertificate, "testVATTaxpayerCertificate", mock.Anything).Return("testVATTaxpayerCertificate").Once()
		s.s3Storage.On("GetPresignedURL", enum.OfficialNewspaperAD, "testOfficialNewspaperAD", mock.Anything).Return("testOfficialNewspaperAD").Once()
		s.s3Storage.On("GetPresignedURL", enum.LogoPic, "testLogo", mock.Anything).Return("testLogo").Once()
		s.addressService.On("GetAddresses", mock.Anything).Return([]addressdto.AddressResponse{}).Once()
		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return(contactInformation).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, uint(1)).Return(contactType, true).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, uint(2)).Return(nilContactType, false).Once()
		s.corporationRepository.On("FindCorporationSignatories", s.db, mock.Anything).Return(signatories).Once()

		request := corporationdto.CorporationDetailsRequest{
			CorporationID: 1,
			Status:        enum.CorpStatusApproved,
		}
		response := s.corporationService.GetCorporationDetails(request)

		s.Equal(response.Name, corporation.Name)
		s.Equal(response.Logo, corporation.Logo)
		s.Equal(response.RegistrationNumber, corporation.RegistrationNumber)
		s.Equal(response.NationalID, corporation.NationalID)
		s.Equal(response.IBAN, corporation.IBAN)
		s.Equal(response.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.s3Storage.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetContactInfo() {
	s.Run("success - Contact info fetched", func() {
		contactInformation := []*entity.ContactInformation{
			{
				TypeID: 1,
				Value:  "testContactValue",
			},
			{
				TypeID: 2,
				Value:  "testContactValue2",
			},
		}
		contactType := &entity.ContactType{
			Name: "testContactType",
		}
		var nilContactType *entity.ContactType = nil

		s.corporationRepository.On("FindContactInformation", s.db, mock.Anything).Return(contactInformation).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, uint(1)).Return(contactType, true).Once()
		s.corporationRepository.On("FindContactTypeByID", s.db, uint(2)).Return(nilContactType, false).Once()

		response := s.corporationService.getContactInfo(1)

		s.Equal(response[0].ContactType.Name, contactType.Name)
		s.Equal(response[0].Value, contactInformation[0].Value)

		s.corporationRepository.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestGetCorporationSignatories() {
	s.Run("success - Corporation signatories fetched", func() {
		signatories := []*entity.Signatory{
			{
				Name: "testSignatory",
			},
		}

		s.corporationRepository.On("FindCorporationSignatories", s.db, mock.Anything).Return(signatories).Once()

		response := s.corporationService.getCorporationSignatories(1)

		s.Equal(response[0].Name, signatories[0].Name)
	})
}

func (s *CorporationServiceTestSuite) TestGetContactTypes() {
	s.Run("success - Contact types fetched", func() {
		contactTypes := []*entity.ContactType{
			{
				Name: "testContactType",
			},
		}

		s.corporationRepository.On("FindContactTypes", s.db).Return(contactTypes).Once()

		response := s.corporationService.GetContactTypes()

		s.Equal(response[0].Name, contactTypes[0].Name)
	})
}

func (s *CorporationServiceTestSuite) TestAddAddress() {
	s.Run("success - Address added", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		addressResponse := addressdto.AddressResponse{
			ID:            1,
			Province:      "testProvince",
			City:          "testCity",
			StreetAddress: "testStreetAddress",
			PostalCode:    "testPostalCode",
			HouseNumber:   "testHouseNumber",
			Unit:          1,
		}
		corporationStaff := &entity.CorporationStaff{}
		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.addressService.On("CreateAddress", mock.Anything).Return(addressResponse).Once()

		request := corporationdto.AddCorporationAddressRequest{
			CorporationID:     1,
			ApplicantID:       1,
			CorporationStatus: enum.CorpStatusApproved,
			Addresses: []addressdto.CreateAddressRequest{
				{
					ProvinceID:    1,
					CityID:        1,
					StreetAddress: "testStreetAddress",
					PostalCode:    "testPostalCode",
					HouseNumber:   "testHouseNumber",
					Unit:          1,
					OwnerID:       1,
				},
			},
		}

		s.corporationService.AddAddress(request)
		s.Equal(addressResponse.StreetAddress, request.Addresses[0].StreetAddress)
		s.Equal(addressResponse.PostalCode, request.Addresses[0].PostalCode)
		s.Equal(addressResponse.HouseNumber, request.Addresses[0].HouseNumber)
		s.Equal(addressResponse.Unit, request.Addresses[0].Unit)

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.AddCorporationAddressRequest{
			CorporationID:     1,
			CorporationStatus: enum.CorpStatusApproved,
			ApplicantID:       1,
		}

		s.Panics(func() {
			s.corporationService.AddAddress(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestDeleteAddress() {
	s.Run("success - Address deleted", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.addressService.On("DeleteAddress", uint(1)).Return(nil).Once()

		request := corporationdto.DeleteAddressRequest{
			CorporationID:     1,
			UserID:            1,
			CorporationStatus: enum.CorpStatusApproved,
			AddressID:         1,
		}
		s.corporationService.DeleteAddress(request)

		s.addressService.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.DeleteAddressRequest{
			CorporationID:     1,
			UserID:            1,
			CorporationStatus: enum.CorpStatusApproved,
			AddressID:         1,
		}

		s.Panics(func() {
			s.corporationService.DeleteAddress(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func (s *CorporationServiceTestSuite) TestChangeLogo() {
	s.Run("success - Logo changed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
			Logo:   "testLogo",
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.LogoPic, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.LogoPic, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, mock.Anything).Return(nil).Once()

		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   1,
			Logo:          &multipart.FileHeader{Filename: "testLogo"},
		}
		s.corporationService.ChangeLogo(request)

		s.s3Storage.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - User not active", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
		}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(false).Once()

		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   1,
			Logo:          &multipart.FileHeader{Filename: "testLogo"},
		}

		s.Panics(func() {
			s.corporationService.ChangeLogo(request)
		})

		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Delete object failed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
			Logo:   "testLogo",
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.LogoPic, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.LogoPic, mock.Anything).Return(errors.New("delete object failed")).Once()

		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   1,
			Logo:          &multipart.FileHeader{Filename: "testLogo"},
		}
		s.Panics(func() {
			s.corporationService.ChangeLogo(request)
		})

		s.s3Storage.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
	s.Run("error - Update corporation failed", func() {
		corporation := &entity.Corporation{
			Status: enum.CorpStatusApproved,
			Logo:   "testLogo",
		}
		corporationStaff := &entity.CorporationStaff{}

		s.corporationRepository.On("FindCorporationByID", s.db, uint(1)).Return(corporation, true).Once()
		s.userService.On("IsUserActive", uint(1)).Return(true).Once()
		s.corporationRepository.On("FindCorporationStaff", s.db, uint(1), uint(1)).Return(corporationStaff, true).Once()
		s.s3Storage.On("UploadObject", enum.LogoPic, mock.Anything, mock.Anything).Return(nil).Once()
		s.s3Storage.On("DeleteObject", enum.LogoPic, mock.Anything).Return(nil).Once()
		s.corporationRepository.On("UpdateCorporation", s.db, mock.Anything).Return(errors.New("update corporation failed")).Once()

		request := corporationdto.ChangeLogoRequest{
			CorporationID: 1,
			ApplicantID:   1,
			Logo:          &multipart.FileHeader{Filename: "testLogo"},
		}
		s.Panics(func() {
			s.corporationService.ChangeLogo(request)
		})

		s.s3Storage.AssertExpectations(s.T())
		s.corporationRepository.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
	})
}

func TestCorporationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CorporationServiceTestSuite))
}
