package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	userService           service.UserService
	addressService        service.AddressService
	s3Storage             s3.S3Storage
	corporationRepository repository.CorporationRepository
	db                    database.Database
}

func NewCorporationService(
	constants *bootstrap.Constants,
	userService service.UserService,
	addressService service.AddressService,
	s3Storage s3.S3Storage,
	corporationRepository repository.CorporationRepository,
	db database.Database,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		userService:           userService,
		addressService:        addressService,
		s3Storage:             s3Storage,
		corporationRepository: corporationRepository,
		db:                    db,
	}
}

// change this logic later if u can
func (corporationService *CorporationService) GetCorporationByID(corporationID uint) *entity.Corporation {
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	return corporation
}

func (corporationService *CorporationService) CheckApplicantAccess(corporationID, applicantID uint) {
	_, exist := corporationService.corporationRepository.FindCorporationStaff(corporationService.db, applicantID, corporationID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.Corporation,
		}
		panic(forbiddenError)
	}
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationDetailsResponse {
	corporationService.userService.GetUserCredential(registerInfo.ApplicantID)
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	_, exist := corporationService.corporationRepository.FindCorporationByName(corporationService.db, registerInfo.Name, activeStatus)
	if exist {
		conflictErrors.Add(corporationService.constants.Field.Name, corporationService.constants.Tag.AlreadyExist)
	}
	_, exist = corporationService.corporationRepository.FindCorporationByNationalID(corporationService.db, registerInfo.NationalID, activeStatus)
	if exist {
		conflictErrors.Add(corporationService.constants.Field.NationalID, corporationService.constants.Tag.AlreadyExist)
	}
	_, exist = corporationService.corporationRepository.FindCorporationByRegistrationNumber(corporationService.db, registerInfo.RegistrationNumber, activeStatus)
	if exist {
		conflictErrors.Add(corporationService.constants.Field.RegistrationNumber, corporationService.constants.Tag.AlreadyExist)
	}
	if registerInfo.IBAN != "" {
		_, exist = corporationService.corporationRepository.FindCorporationByIBAN(corporationService.db, registerInfo.IBAN, activeStatus)
		if exist {
			conflictErrors.Add(corporationService.constants.Field.IBAN, corporationService.constants.Tag.AlreadyExist)
		}
	}
	if len(conflictErrors.Errors) > 0 {
		panic(conflictErrors)
	}

	corporation := &entity.Corporation{
		Name:               registerInfo.Name,
		RegistrationNumber: registerInfo.RegistrationNumber,
		NationalID:         registerInfo.NationalID,
		IBAN:               registerInfo.IBAN,
		Status:             enum.CorpStatusAwaitingApproval,
	}

	err := corporationService.corporationRepository.CreateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}

	staff := &entity.CorporationStaff{
		StaffID:       registerInfo.ApplicantID,
		CorporationID: corporation.ID,
		StaffType:     enum.StaffTypeManager,
	}

	err = corporationService.corporationRepository.CreateCorporationStaff(corporationService.db, staff)
	if err != nil {
		panic(err)
	}

	for _, signatory := range registerInfo.Signatories {
		_, exist = corporationService.corporationRepository.FindCorporationSignatoryByNationalID(corporationService.db, corporation.ID, signatory.NationalCardNumber, signatory.Position)
		signatoryEntity := &entity.Signatory{
			CorporationID:      corporation.ID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
		// errors could be ignored instead of panicking
		err := corporationService.corporationRepository.CreateSignatory(corporationService.db, signatoryEntity)
		if err != nil {
			panic(err)
		}
	}

	return corporationdto.CorporationDetailsResponse{ID: corporation.ID, Name: corporation.Name}
}

func (corporationService *CorporationService) AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) {
	corporation := corporationService.GetCorporationByID(requestInfo.CorporationID)
	corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.ApplicantID)
	if requestInfo.VATTaxpayerCertificate != nil {
		taxPayerPath := corporationService.constants.S3BucketPath.GetVATTaxpayerCertificatePath(corporation.ID, requestInfo.VATTaxpayerCertificate.Filename)
		corporationService.s3Storage.UploadObject(enum.VATTaxpayerCertificate, taxPayerPath, requestInfo.VATTaxpayerCertificate)
		corporation.VATTaxpayerCertificate = taxPayerPath
	}
	if requestInfo.OfficialNewspaperAD != nil {
		newspaperADPath := corporationService.constants.S3BucketPath.GetOfficialNewspaperADPath(corporation.ID, requestInfo.OfficialNewspaperAD.Filename)
		corporationService.s3Storage.UploadObject(enum.OfficialNewspaperAD, newspaperADPath, requestInfo.OfficialNewspaperAD)
		corporation.OfficialNewspaperAD = newspaperADPath
	}
	err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) {
	corporationService.GetCorporationByID(contactInfo.CorporationID)
	corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID)
	for _, contact := range contactInfo.ContactInformation {
		_, exist := corporationService.corporationRepository.FindContactInformationTypeByID(corporationService.db, contact.ContactTypeID)
		if !exist {
			continue
		}
		_, exist = corporationService.corporationRepository.FindContactInformationTypeValue(corporationService.db, contact.ContactTypeID, contact.ContactValue)
		if exist {
			continue
		}
		contact := &entity.ContactInformation{
			CorporationID: contactInfo.CorporationID,
			TypeID:        contact.ContactTypeID,
			Value:         contact.ContactValue,
		}
		err := corporationService.corporationRepository.CreateContactInformation(corporationService.db, contact)
		if err != nil {
			panic(err)
		}
	}
}

func (corporationService *CorporationService) AddAddress(addressInfo corporationdto.AddCorporationAddressRequest) {
	corporationService.GetCorporationByID(addressInfo.CorporationID)
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.ApplicantID)
	for _, address := range addressInfo.Addresses {
		corporationService.addressService.CreateAddress(address)
	}
}

func (corporationService *CorporationService) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) {
	corporationService.GetCorporationByID(addressInfo.CorporationID)
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.UserID)
	corporationService.addressService.DeleteAddress(addressInfo.AddressID)
}
