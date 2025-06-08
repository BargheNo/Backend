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

func (corporationService *CorporationService) GetCorporationStatuses() []corporationdto.GetStatusesResponse {
	statuses := enum.GetAllCorporationStatuses()
	response := make([]corporationdto.GetStatusesResponse, len(statuses))
	for i, status := range statuses {
		response[i] = corporationdto.GetStatusesResponse{
			ID:     uint(status),
			Status: status.String(),
		}
	}
	return response
}

func (corporationService *CorporationService) getCorporationByIDAndStatus(corporationID uint, status enum.CorporationStatus) (*entity.Corporation, error) {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return nil, err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return nil, notFoundError
	}
	if corporation.Status != status {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.Corporation,
		}
		return nil, &forbiddenError
	}
	return corporation, nil
}

func (corporationService *CorporationService) DoesCorporationExist(corporationID uint) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return notFoundError
	}
	return nil
}

func (corporationService *CorporationService) GetCorporationCredentials(corporationID uint) (corporationdto.CorporationCredentialResponse, error) {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return corporationdto.CorporationCredentialResponse{}, notFoundError
	}
	ownerInfo := addressdto.GetOwnerAddressesRequest{
		OwnerID:   corporation.ID,
		OwnerType: corporationService.constants.AddressOwners.Corporation,
	}
	addresses, err := corporationService.addressService.GetAddresses(ownerInfo)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	contactInfo, err := corporationService.getContactInfo(corporation.ID)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	return corporationdto.CorporationCredentialResponse{
		ID:          corporation.ID,
		Name:        corporation.Name,
		ContactInfo: contactInfo,
		Addresses:   addresses,
	}, nil
}

func (corporationService *CorporationService) ISCorporationApproved(corporationID uint) (bool, error) {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return false, err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return false, notFoundError
	}
	isApproved := corporation.Status == enum.CorpStatusApproved
	return isApproved, nil
}

// TODO: need to change check applicant access to only one function and if it was failed got 404 not found
func (corporationService *CorporationService) CheckApplicantAccess(corporationID, applicantID uint) error {
	staff, err := corporationService.corporationRepository.FindCorporationStaff(corporationService.db, applicantID, corporationID)
	if err != nil {
		return err
	}
	if staff == nil {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.Corporation,
		}
		return &forbiddenError
	}
	return nil
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) (corporationdto.CorporationCredentialResponse, error) {
	ok, err := corporationService.userService.IsUserActive(registerInfo.ApplicantID)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return corporationdto.CorporationCredentialResponse{}, &forbiddenError
	}
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	corporation, err := corporationService.corporationRepository.FindCorporationByName(corporationService.db, registerInfo.Name, activeStatus)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.Name, corporationService.constants.Tag.AlreadyExist)
	}
	corporation, err = corporationService.corporationRepository.FindCorporationByNationalID(corporationService.db, registerInfo.NationalID, activeStatus)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.NationalID, corporationService.constants.Tag.AlreadyExist)
	}
	corporation, err = corporationService.corporationRepository.FindCorporationByRegistrationNumber(corporationService.db, registerInfo.RegistrationNumber, activeStatus)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.RegistrationNumber, corporationService.constants.Tag.AlreadyExist)
	}
	if registerInfo.IBAN != "" {
		corporation, err = corporationService.corporationRepository.FindCorporationByIBAN(corporationService.db, registerInfo.IBAN, activeStatus)
		if err != nil {
			return corporationdto.CorporationCredentialResponse{}, err
		}
		if corporation != nil {
			conflictErrors.Add(corporationService.constants.Field.IBAN, corporationService.constants.Tag.AlreadyExist)
		}
	}
	if len(conflictErrors.Errors) > 0 {
		return corporationdto.CorporationCredentialResponse{}, conflictErrors
	}

	corporation = &entity.Corporation{
		Name:               registerInfo.Name,
		RegistrationNumber: registerInfo.RegistrationNumber,
		NationalID:         registerInfo.NationalID,
		IBAN:               registerInfo.IBAN,
		Status:             enum.CorpStatusAwaitingApproval,
	}

	err = corporationService.corporationRepository.CreateCorporation(corporationService.db, corporation)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	staff := &entity.CorporationStaff{
		StaffID:       registerInfo.ApplicantID,
		CorporationID: corporation.ID,
		StaffType:     enum.StaffTypeManager,
	}

	err = corporationService.corporationRepository.CreateCorporationStaff(corporationService.db, staff)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	for _, signatory := range registerInfo.Signatories {
		signatory, err := corporationService.corporationRepository.FindCorporationSignatoryByNationalID(corporationService.db, corporation.ID, signatory.NationalCardNumber, signatory.Position)
		if err != nil {
			return corporationdto.CorporationCredentialResponse{}, err
		}
		if signatory != nil {
			continue
		}
		signatoryEntity := &entity.Signatory{
			CorporationID:      corporation.ID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
		err = corporationService.corporationRepository.CreateSignatory(corporationService.db, signatoryEntity)
		if err != nil {
			return corporationdto.CorporationCredentialResponse{}, err
		}
	}

	return corporationdto.CorporationCredentialResponse{ID: corporation.ID, Name: corporation.Name}, nil
}

func (corporationService *CorporationService) replaceSignatories(corporationID uint, Signatories []corporationdto.Signatory) error {
	err := corporationService.corporationRepository.DeleteCorporationSignatories(corporationService.db, corporationID)
	if err != nil {
		return err
	}
	for _, signatory := range Signatories {
		signatory, err := corporationService.corporationRepository.FindCorporationSignatoryByNationalID(corporationService.db, corporationID, signatory.NationalCardNumber, signatory.Position)
		if err != nil {
			return err
		}
		if signatory != nil {
			continue
		}
		signatoryEntity := &entity.Signatory{
			CorporationID:      corporationID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
		err = corporationService.corporationRepository.CreateSignatory(corporationService.db, signatoryEntity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (corporationService *CorporationService) UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(updateRegisterInfo.CorporationID, enum.CorpStatusAwaitingApproval)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(updateRegisterInfo.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
	}
	corporationService.CheckApplicantAccess(updateRegisterInfo.CorporationID, updateRegisterInfo.ApplicantID)

	corporationService.checkCorporationConflicts(corporation, updateRegisterInfo.Name, updateRegisterInfo.NationalID, updateRegisterInfo.RegistrationNumber, updateRegisterInfo.IBAN)

	err = corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		return err
	}

	err = corporationService.replaceSignatories(updateRegisterInfo.CorporationID, updateRegisterInfo.Signatories)
	if err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) checkCorporationConflicts(corporation *entity.Corporation, name, nationalID, registrationNumber, iban *string) error {
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	if name != nil {
		corporation, err := corporationService.corporationRepository.FindCorporationByName(corporationService.db, *name, activeStatus)
		if err != nil {
			return err
		}
		if corporation != nil {
			conflictErrors.Add(corporationService.constants.Field.Name, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.Name = *name
	}

	if nationalID != nil {
		corporation, err := corporationService.corporationRepository.FindCorporationByNationalID(corporationService.db, *nationalID, activeStatus)
		if err != nil {
			return err
		}
		if corporation != nil {
			conflictErrors.Add(corporationService.constants.Field.NationalID, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.NationalID = *nationalID
	}

	if registrationNumber != nil {
		corporation, err := corporationService.corporationRepository.FindCorporationByRegistrationNumber(corporationService.db, *registrationNumber, activeStatus)
		if err != nil {
			return err
		}
		if corporation != nil {
			conflictErrors.Add(corporationService.constants.Field.RegistrationNumber, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.RegistrationNumber = *registrationNumber
	}

	if iban != nil {
		corporation, err := corporationService.corporationRepository.FindCorporationByIBAN(corporationService.db, *iban, activeStatus)
		if err != nil {
			return err
		}
		if corporation != nil {
			conflictErrors.Add(corporationService.constants.Field.IBAN, corporationService.constants.Tag.AlreadyExist)
		}
		corporation.IBAN = *iban
	}

	if len(conflictErrors.Errors) > 0 {
		return conflictErrors
	}
	return nil
}

func (corporationService *CorporationService) AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(requestInfo.CorporationID, enum.CorpStatusAwaitingApproval)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(requestInfo.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
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
				return err
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
				return err
			}
		}
	}
	err = corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) error {
	_, err := corporationService.getCorporationByIDAndStatus(contactInfo.CorporationID, contactInfo.CorporationStatus)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(contactInfo.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
	}
	corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID)

	for _, contact := range contactInfo.ContactInformation {
		contactType, err := corporationService.corporationRepository.FindContactInformationTypeByID(corporationService.db, contact.ContactTypeID)
		if err != nil {
			return err
		}
		if contactType == nil {
			continue
		}
		contactInfo, err := corporationService.corporationRepository.FindContactInformationTypeValue(corporationService.db, contact.ContactTypeID, contact.ContactValue)
		if err != nil {
			return err
		}
		if contactInfo != nil {
			continue
		}
		contact := &entity.ContactInformation{
			CorporationID: contactInfo.CorporationID,
			TypeID:        contact.ContactTypeID,
			Value:         contact.ContactValue,
		}
		err = corporationService.corporationRepository.CreateContactInformation(corporationService.db, contact)
		if err != nil {
			return err
		}
	}
	return nil
}

func (corporationService *CorporationService) DeleteContactInfo(contactInfo corporationdto.DeleteContactInformationRequest) error {
	_, err := corporationService.getCorporationByIDAndStatus(contactInfo.CorporationID, contactInfo.CorporationStatus)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(contactInfo.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
	}
	corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID)

	contact, err := corporationService.corporationRepository.FindContactInformationByID(corporationService.db, contactInfo.ContactID)
	if err != nil {
		return err
	}
	if contact == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.ContactInformation}
		return &notFoundError
	}

	err = corporationService.corporationRepository.DeleteContactInfo(corporationService.db, contact)
	if err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) getPrivateCorporationDetails(corporation *entity.Corporation) (corporationdto.CorporationPrivateInfoResponse, error) {
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
	addresses, err := corporationService.addressService.GetAddresses(ownerInfo)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	contactInfo, err := corporationService.getContactInfo(corporation.ID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	signatories, err := corporationService.getCorporationSignatories(corporation.ID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

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
	}, nil
}

func (corporationService *CorporationService) GetCorporationDetails(requestInfo corporationdto.CorporationDetailsRequest) (corporationdto.CorporationPrivateInfoResponse, error) {
	err := corporationService.userService.DoesUserExist(requestInfo.UserID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	corporation, err := corporationService.getCorporationByIDAndStatus(requestInfo.CorporationID, requestInfo.Status)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	err = corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.UserID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	details, err := corporationService.getPrivateCorporationDetails(corporation)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	return details, nil
}

func (corporationService *CorporationService) getContactInfo(corporationID uint) ([]corporationdto.ContactInformationResponse, error) {
	contactInfo, err := corporationService.corporationRepository.FindContactInformation(corporationService.db, corporationID)
	if err != nil {
		return nil, err
	}
	response := make([]corporationdto.ContactInformationResponse, len(contactInfo))
	for i, contact := range contactInfo {
		contactType, err := corporationService.corporationRepository.FindContactTypeByID(corporationService.db, contact.TypeID)
		if err != nil {
			return nil, err
		}
		if contactType == nil {
			continue
		}
		response[i] = corporationdto.ContactInformationResponse{
			ID:          contact.ID,
			ContactType: corporationdto.ContactTypeResponse{ID: contactType.ID, Name: contactType.Name},
			Value:       contact.Value,
		}
	}
	return response, nil
}

func (corporationService *CorporationService) getCorporationSignatories(corporationID uint) ([]corporationdto.SignatoryResponse, error) {
	signatories, err := corporationService.corporationRepository.FindCorporationSignatories(corporationService.db, corporationID)
	if err != nil {
		return nil, err
	}
	response := make([]corporationdto.SignatoryResponse, len(signatories))
	for i, signatory := range signatories {
		response[i] = corporationdto.SignatoryResponse{
			ID:                 signatory.ID,
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
	}
	return response, nil
}

func (corporationService *CorporationService) GetContactTypes() ([]corporationdto.ContactTypeResponse, error) {
	types, err := corporationService.corporationRepository.FindContactTypes(corporationService.db)
	if err != nil {
		return nil, err
	}
	contactTypes := make([]corporationdto.ContactTypeResponse, len(types))
	for i, contactType := range types {
		contactTypes[i] = corporationdto.ContactTypeResponse{
			ID:   contactType.ID,
			Name: contactType.Name,
		}
	}
	return contactTypes, nil
}

func (corporationService *CorporationService) AddAddress(addressInfo corporationdto.AddCorporationAddressRequest) error {
	_, err := corporationService.getCorporationByIDAndStatus(addressInfo.CorporationID, addressInfo.CorporationStatus)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(addressInfo.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
	}
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.ApplicantID)

	for _, address := range addressInfo.Addresses {
		corporationService.addressService.CreateAddress(address)
	}
	return nil
}

func (corporationService *CorporationService) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) error {
	_, err := corporationService.getCorporationByIDAndStatus(addressInfo.CorporationID, addressInfo.CorporationStatus)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(addressInfo.UserID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
	}
	corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.UserID)

	corporationService.addressService.DeleteAddress(addressInfo.AddressID)
	return nil
}

func (corporationService *CorporationService) ChangeLogo(changeLogoRequest corporationdto.ChangeLogoRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(changeLogoRequest.CorporationID, enum.CorpStatusApproved)
	if err != nil {
		return err
	}

	ok, err := corporationService.userService.IsUserActive(changeLogoRequest.ApplicantID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: "",
		}
		return &forbiddenError
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
				return err
			}
		}
	}
	err = corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) GetCorporationByIDAndStatus(corporationID uint, status enum.CorporationStatus) (*entity.Corporation, error) {
	return corporationService.getCorporationByIDAndStatus(corporationID, status)
}

func (corporationService *CorporationService) GetUserCorporations(userID uint) ([]corporationdto.CorporationCredentialResponse, error) {
	corporations, err := corporationService.corporationRepository.FindUserCorporations(corporationService.db, userID)
	if err != nil {
		return nil, err
	}

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		credentials, err := corporationService.GetCorporationCredentials(corporation.ID)
		if err != nil {
			return nil, err
		}
		response[i] = credentials
	}
	return response, nil
}

func (corporationService *CorporationService) GetAvailableCorporations() ([]corporationdto.CorporationCredentialResponse, error) {
	allowedStatuses := []enum.CorporationStatus{enum.CorpStatusApproved}
	corporations, err := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses)
	if err != nil {
		return nil, err
	}

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		credentials, err := corporationService.GetCorporationCredentials(corporation.ID)
		if err != nil {
			return nil, err
		}
		response[i] = credentials
	}
	return response, nil
}

func (corporationService *CorporationService) GetCorporationsByAdmin(listInfo corporationdto.GetCorporationsByAdminRequest) ([]corporationdto.CorporationCredentialResponse, error) {
	allowedStatuses := corporationService.mapStatusIDToAllowedStatuses(listInfo.Status)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	corporations, err := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses, sortingModifier, paginationModifier)
	if err != nil {
		return nil, err
	}

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		credentials, err := corporationService.GetCorporationCredentials(corporation.ID)
		if err != nil {
			return nil, err
		}
		response[i] = credentials
	}
	return response, nil
}

func (corporationService *CorporationService) GetCorporationByAdmin(corporationID uint) (corporationdto.CorporationPrivateInfoResponse, error) {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return corporationdto.CorporationPrivateInfoResponse{}, &notFoundError
	}
	details, err := corporationService.getPrivateCorporationDetails(corporation)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	return details, nil
}

func (corporationService *CorporationService) GetReviewActions() ([]corporationdto.GetStatusesResponse, error) {
	actions := enum.GetAllReviewActions()
	response := make([]corporationdto.GetStatusesResponse, len(actions))
	for i, action := range actions {
		response[i] = corporationdto.GetStatusesResponse{
			ID:     uint(action),
			Status: action.String(),
		}
	}
	return response, nil
}

func (corporationService *CorporationService) GetCorporationReviewsByAdmin(corporationID uint) ([]corporationdto.GetAdminCorporationReview, error) {
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	reviews, err := corporationService.corporationRepository.FindCorporationReviews(corporationService.db, corporationID, sortingModifier)
	if err != nil {
		return nil, err
	}

	response := make([]corporationdto.GetAdminCorporationReview, len(reviews))
	for i, review := range reviews {
		operator, err := corporationService.userService.GetUserCredential(review.ReviewerID)
		if err != nil {
			return nil, err
		}
		response[i] = corporationdto.GetAdminCorporationReview{
			Reviewer: operator,
			Action:   review.Action.String(),
			Reason:   review.Reason,
			Notes:    review.Notes,
		}
	}
	return response, nil
}

func (corporationService *CorporationService) ApproveCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, request.CorporationID)
	if err != nil {
		return err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return &notFoundError
	}

	var conflictErrors exception.ConflictErrors
	if corporation.Status == enum.CorpStatusApproved {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyAccepted)
		return &conflictErrors
	} else if corporation.Status == enum.CorpStatusRejected {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyRejected)
		return &conflictErrors
	} else if corporation.Status != enum.CorpStatusAwaitingApproval {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.ForbiddenStatus)
		return &conflictErrors
	}

	review := &entity.CorporationReview{
		CorporationID: request.CorporationID,
		ReviewerID:    request.ReviewerID,
		Action:        enum.ReviewActionApproved,
		Reason:        request.Reason,
		Notes:         request.Notes,
	}
	if err := corporationService.corporationRepository.CreateReview(corporationService.db, review); err != nil {
		return err
	}

	corporation.Status = enum.CorpStatusApproved
	if err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation); err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) RejectCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, request.CorporationID)
	if err != nil {
		return err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return &notFoundError
	}

	if enum.ReviewAction(request.ActionID) == enum.ReviewActionApproved {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.CorporationReview,
		}
		return &forbiddenError
	}

	var conflictErrors exception.ConflictErrors
	if corporation.Status == enum.CorpStatusApproved {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyAccepted)
		return &conflictErrors
	} else if corporation.Status == enum.CorpStatusRejected {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyRejected)
		return &conflictErrors
	} else if corporation.Status != enum.CorpStatusAwaitingApproval {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.ForbiddenStatus)
		return &conflictErrors
	}

	review := &entity.CorporationReview{
		CorporationID: request.CorporationID,
		ReviewerID:    request.ReviewerID,
		Action:        enum.ReviewAction(request.ActionID),
		Reason:        request.Reason,
		Notes:         request.Notes,
	}
	if err := corporationService.corporationRepository.CreateReview(corporationService.db, review); err != nil {
		return err
	}

	var corpStatus enum.CorporationStatus
	if enum.ReviewAction(request.ActionID) == enum.ReviewActionSuspended {
		corpStatus = enum.CorpStatusSuspend
	} else {
		corpStatus = enum.CorpStatusRejected
	}

	corporation.Status = corpStatus
	if err := corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation); err != nil {
		return err
	}
	return nil
}
