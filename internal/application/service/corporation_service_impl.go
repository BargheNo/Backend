package serviceimpl

import (
	"time"

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

func (corporationService *CorporationService) mapStatusIDToAllowedStatuses(statusID uint) []enum.CorporationStatus {
	status := enum.CorporationStatus(statusID)

	allowedStatuses := enum.GetAllCorporationStatuses()

	for _, allowedStatus := range allowedStatuses {
		if status == allowedStatus {
			if status == enum.CorpStatusAll {
				return allowedStatuses
			}
			return []enum.CorporationStatus{status}
		}
	}
	return allowedStatuses
}

func (corporationService *CorporationService) GetCorporationStatuses() []corporationdto.GetCorporationStatusesResponse {
	statuses := enum.GetAllCorporationStatuses()
	response := make([]corporationdto.GetCorporationStatusesResponse, len(statuses))
	for i, status := range statuses {
		response[i] = corporationdto.GetCorporationStatusesResponse{
			ID:     uint(status),
			Status: status.String(),
		}
	}
	return response
}

func (corporationService *CorporationService) getCorporationByIDAndStatus(corporationID uint, status enum.CorporationStatus) *entity.Corporation {
	corporation, exist := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	if corporation.Status != status {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.Corporation,
		}
		panic(forbiddenError)
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

func (corporationService *CorporationService) GetCorporationCredentials(corporationID uint) corporationdto.CorporationCredentialResponse {
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
	contactInfo := corporationService.getContactInfo(corporation.ID)
	return corporationdto.CorporationCredentialResponse{
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

// TODO: need to change check applicant access to only one function and if it was failed got 404 not found
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

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationCredentialResponse {
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

	return corporationdto.CorporationCredentialResponse{ID: corporation.ID, Name: corporation.Name}
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
	corporation := corporationService.getCorporationByIDAndStatus(updateRegisterInfo.CorporationID, enum.CorpStatusAwaitingApproval)

	exist := corporationService.userService.IsUserActive(updateRegisterInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(updateRegisterInfo.CorporationID, updateRegisterInfo.ApplicantID)

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
	corporation := corporationService.getCorporationByIDAndStatus(requestInfo.CorporationID, enum.CorpStatusAwaitingApproval)

	exist := corporationService.userService.IsUserActive(requestInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.ApplicantID)

	prevVatTaxPayerPath := corporation.VATTaxpayerCertificate
	prevOfficialNewspaperPath := corporation.OfficialNewspaperAD

	if requestInfo.VATTaxpayerCertificate != nil {
		taxPayerPath := corporationService.constants.S3BucketPath.GetVATTaxpayerCertificatePath(corporation.ID, requestInfo.VATTaxpayerCertificate.Filename)
		corporationService.s3Storage.UploadObject(enum.VATTaxpayerCertificate, taxPayerPath, requestInfo.VATTaxpayerCertificate)
		corporation.VATTaxpayerCertificate = taxPayerPath
		if prevVatTaxPayerPath != "" {
			err := corporationService.s3Storage.DeleteObject(enum.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate)
			if err != nil {
				panic(err)
			}
		}
	}

	if requestInfo.OfficialNewspaperAD != nil {
		newspaperADPath := corporationService.constants.S3BucketPath.GetOfficialNewspaperADPath(corporation.ID, requestInfo.OfficialNewspaperAD.Filename)
		corporationService.s3Storage.UploadObject(enum.OfficialNewspaperAD, newspaperADPath, requestInfo.OfficialNewspaperAD)
		corporation.OfficialNewspaperAD = newspaperADPath
		if prevOfficialNewspaperPath != "" {
			err := corporationService.s3Storage.DeleteObject(enum.OfficialNewspaperAD, corporation.OfficialNewspaperAD)
			if err != nil {
				panic(err)
			}
		}
	}
	err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) {
	corporationService.getCorporationByIDAndStatus(contactInfo.CorporationID, contactInfo.CorporationStatus)

	exist := corporationService.userService.IsUserActive(contactInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
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

func (corporationService *CorporationService) DeleteContactInfo(contactInfo corporationdto.DeleteContactInformationRequest) {
	corporationService.getCorporationByIDAndStatus(contactInfo.CorporationID, contactInfo.CorporationStatus)

	exist := corporationService.userService.IsUserActive(contactInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID)

	contact, exist := corporationService.corporationRepository.FindContactInformationByID(corporationService.db, contactInfo.ContactID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.ContactInformation}
		panic(notFoundError)
	}

	if err := corporationService.corporationRepository.DeleteContactInfo(corporationService.db, contact); err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) GetCorporationDetails(requestInfo corporationdto.CorporationDetailsRequest) corporationdto.CorporationPrivateInfoResponse {
	corporationService.userService.DoesUserExist(requestInfo.UserID)
	corporation := corporationService.getCorporationByIDAndStatus(requestInfo.CorporationID, requestInfo.Status)
	corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.UserID)

	vatTaxPayer := ""
	if corporation.VATTaxpayerCertificate != "" {
		vatTaxPayer = corporationService.s3Storage.GetPresignedURL(enum.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate, 8*time.Hour)
	}

	officialNewspaperAD := ""
	if corporation.OfficialNewspaperAD != "" {
		officialNewspaperAD = corporationService.s3Storage.GetPresignedURL(enum.OfficialNewspaperAD, corporation.OfficialNewspaperAD, 8*time.Hour)
	}

	logo := ""
	if corporation.Logo != "" {
		logo = corporationService.s3Storage.GetPresignedURL(enum.LogoPic, corporation.Logo, 8*time.Hour)
	}

	ownerInfo := addressdto.GetOwnerAddressesRequest{
		OwnerID:   corporation.ID,
		OwnerType: corporationService.constants.AddressOwners.Corporation,
	}
	addresses := corporationService.addressService.GetAddresses(ownerInfo)

	contactInfo := corporationService.getContactInfo(corporation.ID)

	signatories := corporationService.getCorporationSignatories(requestInfo.CorporationID)

	return corporationdto.CorporationPrivateInfoResponse{
		ID:                     corporation.ID,
		Name:                   corporation.Name,
		Logo:                   logo,
		RegistrationNumber:     corporation.RegistrationNumber,
		NationalID:             corporation.NationalID,
		IBAN:                   corporation.IBAN,
		VATTaxpayerCertificate: vatTaxPayer,
		OfficialNewspaperAD:    officialNewspaperAD,
		Signatories:            signatories,
		ContactInfo:            contactInfo,
		Addresses:              addresses,
	}
}

func (corporationService *CorporationService) getContactInfo(corporationID uint) []corporationdto.ContactInformationResponse {
	contactInfo := corporationService.corporationRepository.FindContactInformation(corporationService.db, corporationID)
	response := make([]corporationdto.ContactInformationResponse, len(contactInfo))
	for i, contact := range contactInfo {
		contactType, exist := corporationService.corporationRepository.FindContactTypeByID(corporationService.db, contact.TypeID)
		if !exist {
			continue
		}
		response[i] = corporationdto.ContactInformationResponse{
			ID:          contact.ID,
			ContactType: corporationdto.ContactTypeResponse{ID: contactType.ID, Name: contactType.Name},
			Value:       contact.Value,
		}
	}
	return response
}

func (corporationService *CorporationService) getCorporationSignatories(corporationID uint) []corporationdto.SignatoryResponse {
	signatories := corporationService.corporationRepository.FindCorporationSignatories(corporationService.db, corporationID)
	response := make([]corporationdto.SignatoryResponse, len(signatories))
	for i, signatory := range signatories {
		response[i] = corporationdto.SignatoryResponse{
			ID:                 signatory.ID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
	}
	return response
}

func (corporationService *CorporationService) GetContactTypes() []corporationdto.ContactTypeResponse {
	types := corporationService.corporationRepository.FindContactTypes(corporationService.db)
	contactTypes := make([]corporationdto.ContactTypeResponse, len(types))
	for i, contactType := range types {
		contactTypes[i] = corporationdto.ContactTypeResponse{
			ID:   contactType.ID,
			Name: contactType.Name,
		}
	}
	return contactTypes
}

func (corporationService *CorporationService) AddAddress(addressInfo corporationdto.AddCorporationAddressRequest) {
	corporationService.getCorporationByIDAndStatus(addressInfo.CorporationID, addressInfo.CorporationStatus)

	exist := corporationService.userService.IsUserActive(addressInfo.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.ApplicantID)

	for _, address := range addressInfo.Addresses {
		corporationService.addressService.CreateAddress(address)
	}
}

func (corporationService *CorporationService) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) {
	corporationService.getCorporationByIDAndStatus(addressInfo.CorporationID, addressInfo.CorporationStatus)

	exist := corporationService.userService.IsUserActive(addressInfo.UserID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.UserID)

	corporationService.addressService.DeleteAddress(addressInfo.AddressID)
}

func (corporationService *CorporationService) ChangeLogo(changeLogoRequest corporationdto.ChangeLogoRequest) {
	corporation := corporationService.getCorporationByIDAndStatus(changeLogoRequest.CorporationID, enum.CorpStatusApproved)

	exist := corporationService.userService.IsUserActive(changeLogoRequest.ApplicantID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		panic(forbiddenError)
	}
	corporationService.CheckApplicantAccess(changeLogoRequest.CorporationID, changeLogoRequest.ApplicantID)

	prevLogoPath := corporation.Logo

	if changeLogoRequest.Logo != nil {
		newLogoPath := corporationService.constants.S3BucketPath.GetCorporationLogoPath(changeLogoRequest.CorporationID, changeLogoRequest.Logo.Filename)
		corporationService.s3Storage.UploadObject(enum.LogoPic, newLogoPath, changeLogoRequest.Logo)
		corporation.Logo = newLogoPath
		if prevLogoPath != "" {
			err := corporationService.s3Storage.DeleteObject(enum.LogoPic, corporation.Logo)
			if err != nil {
				panic(err)
			}
		}
	}
	err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) GetCorporationByIDAndStatus(corporationID uint, status enum.CorporationStatus) *entity.Corporation {
	return corporationService.getCorporationByIDAndStatus(corporationID, status)
}

func (corporationService *CorporationService) GetUserCorporations(userID uint) []corporationdto.CorporationCredentialResponse {
	corporations := corporationService.corporationRepository.FindUserCorporations(corporationService.db, userID)

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		response[i] = corporationService.GetCorporationCredentials(corporation.ID)
	}
	return response
}

func (corporationService *CorporationService) GetAvailableCorporations() []corporationdto.CorporationCredentialResponse {
	allowedStatuses := []enum.CorporationStatus{enum.CorpStatusApproved}
	corporations := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses)

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		response[i] = corporationService.GetCorporationCredentials(corporation.ID)
	}
	return response
}

func (corporationService *CorporationService) GetCorporationsByAdmin(listInfo corporationdto.GetCorporationsByAdminRequest) []corporationdto.CorporationCredentialResponse {
	allowedStatuses := corporationService.mapStatusIDToAllowedStatuses(listInfo.Status)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	corporations := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses, sortingModifier, paginationModifier)

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		response[i] = corporationService.GetCorporationCredentials(corporation.ID)
	}
	return response
}
