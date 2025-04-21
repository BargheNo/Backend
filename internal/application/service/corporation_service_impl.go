package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
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

func (corporationService *CorporationService) getCorporationByID(corporationID uint) *entity.Corporation {
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	return corporation
}

func (corporationService *CorporationService) DoesCorporationExist(corporationID uint) {
	_, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
}

func (corporationService *CorporationService) GetCorporationCredentials(corporationID uint) corporationdto.CorporationDetailsResponse {
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	ownerInfo := addressdto.GetOwnerAddressesRequest{
		OwnerID:   corporation.ID,
		OwnerType: corporationService.constants.AddressOwners.Corporation,
	}
	addresses := corporationService.addressService.GetAddresses(ownerInfo)
	contactInfo := corporationService.GetContactInfo(corporation.ID)
	return corporationdto.CorporationDetailsResponse{
		ID:          corporation.ID,
		Name:        corporation.Name,
		ContactInfo: contactInfo,
		Addresses:   addresses,
	}
}

func (corporationService *CorporationService) ISCorporationApproved(corporationID uint) bool {
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	isApproved := corporation.Status == enum.CorpStatusApproved
	return isApproved
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
	exist := corporationService.userService.IsUserActive(registerInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	_, exist = corporationService.corporationRepository.FindCorporationByName(corporationService.db, registerInfo.Name, activeStatus)
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
		if exist {
			continue
		}
		signatoryEntity := &entity.Signatory{
			CorporationID:      corporation.ID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
		err := corporationService.corporationRepository.CreateSignatory(corporationService.db, signatoryEntity)
		if err != nil {
			panic(err)
		}
	}

	return corporationdto.CorporationDetailsResponse{ID: corporation.ID, Name: corporation.Name}
}

func (corporationService *CorporationService) replaceSignatories(corporationID uint, Signatories []corporationdto.Signatory) {
	err := corporationService.corporationRepository.DeleteCorporationSignatories(corporationService.db, corporationID)
	if err != nil {
		panic(err)
	}
	for _, signatory := range Signatories {
		_, exist := corporationService.corporationRepository.FindCorporationSignatoryByNationalID(corporationService.db, corporationID, signatory.NationalCardNumber, signatory.Position)
		if exist {
			continue
		}
		signatoryEntity := &entity.Signatory{
			CorporationID:      corporationID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
		err := corporationService.corporationRepository.CreateSignatory(corporationService.db, signatoryEntity)
		if err != nil {
			panic(err)
		}
	}
}

func (corporationService *CorporationService) UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest) {
	exist := corporationService.userService.IsUserActive(updateRegisterInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, updateRegisterInfo.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	corporationService.checkCorporationConflicts(corporation, updateRegisterInfo.Name, updateRegisterInfo.NationalID, updateRegisterInfo.RegistrationNumber, updateRegisterInfo.IBAN)

	if err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation); err != nil {
		panic(err)
	}

	corporationService.replaceSignatories(updateRegisterInfo.CorporationID, updateRegisterInfo.Signatories)
}

func (corporationService *CorporationService) checkCorporationConflicts(corporation *entity.Corporation, name, nationalID, registrationNumber, iban *string) {
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	if name != nil {
		_, exist := corporationService.corporationRepository.FindCorporationByName(corporationService.db, *name, activeStatus)
		if exist {
			conflictErrors.Add(corporationService.constants.Field.Name, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.Name = *name
	}

	if nationalID != nil {
		_, exist := corporationService.corporationRepository.FindCorporationByNationalID(corporationService.db, *nationalID, activeStatus)
		if exist {
			conflictErrors.Add(corporationService.constants.Field.NationalID, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.NationalID = *nationalID
	}

	if registrationNumber != nil {
		_, exist := corporationService.corporationRepository.FindCorporationByRegistrationNumber(corporationService.db, *registrationNumber, activeStatus)
		if exist {
			conflictErrors.Add(corporationService.constants.Field.RegistrationNumber, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.RegistrationNumber = *registrationNumber
	}

	if iban != nil {
		_, exist := corporationService.corporationRepository.FindCorporationByIBAN(corporationService.db, *iban, activeStatus)
		if exist {
			conflictErrors.Add(corporationService.constants.Field.IBAN, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.IBAN = *iban
	}

	if len(conflictErrors.Errors) > 0 {
		panic(conflictErrors)
	}
}

func (corporationService *CorporationService) AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) {
	corporation := corporationService.getCorporationByID(requestInfo.CorporationID)
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
	corporationService.DoesCorporationExist(contactInfo.CorporationID)
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
	corporationService.DoesCorporationExist(addressInfo.CorporationID)
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.ApplicantID)
	for _, address := range addressInfo.Addresses {
		corporationService.addressService.CreateAddress(address)
	}
}

func (corporationService *CorporationService) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) {
	corporationService.DoesCorporationExist(addressInfo.CorporationID)
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.UserID)
	corporationService.addressService.DeleteAddress(addressInfo.AddressID)
}

func (corporationService *CorporationService) GetCorporations(requestInfo corporationdto.CorporationListRequest) []corporationdto.CorporationDetailsResponse {
	corporationService.userService.DoesUserExist(requestInfo.UserID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	allowedStatuses := []enum.CorporationStatus{enum.CorpStatusApproved}
	corporations := corporationService.corporationRepository.FindCorporationByStatus(corporationService.db, allowedStatuses, paginationModifier, sortingModifier)
	response := make([]corporationdto.CorporationDetailsResponse, len(corporations))
	for i, corporation := range corporations {
		response[i] = corporationService.GetCorporationCredentials(corporation.ID)
	}

	return response
}

func (corporationService *CorporationService) GetContactInfo(corporationID uint) []corporationdto.ContactInformationResponse {
	corporation := corporationService.getCorporationByID(corporationID)
	contactInfo := corporationService.corporationRepository.FindContactInformation(corporationService.db, corporation.ID)
	response := make([]corporationdto.ContactInformationResponse, len(contactInfo))
	for i, contact := range contactInfo {
		response[i] = corporationdto.ContactInformationResponse{
			ContactTypeID: contact.TypeID,
			ContactValue:  contact.Value,
		}
	}
	return response
}
