package service

import (
	"log"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	userService           usecase.UserService
	addressService        usecase.AddressService
	rbacService           usecase.RBACService
	s3Storage             s3.S3Storage
	corporationRepository postgres.CorporationRepository
	db                    database.Database
}

func NewCorporationService(
	constants *bootstrap.Constants,
	userService usecase.UserService,
	addressService usecase.AddressService,
	rbacService usecase.RBACService,
	s3Storage s3.S3Storage,
	corporationRepository postgres.CorporationRepository,
	db database.Database,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		userService:           userService,
		addressService:        addressService,
		rbacService:           rbacService,
		s3Storage:             s3Storage,
		corporationRepository: corporationRepository,
		db:                    db,
	}
}

func (corporationService *CorporationService) getSortByColumn(requested uint) string {
	allowed := sortby.GetCorporationSortableColumns()
	sortBy := sortby.CorporationSortBy(requested)
	if _, ok := allowed[sortBy]; ok {
		return sortBy.DBColumn()
	}
	return sortby.NewsSortByCreatedAt.DBColumn()
}

func (corporationService *CorporationService) getStaffSortByColumn(requested uint) string {
	allowed := sortby.GetStaffSortableColumns()
	sortBy := sortby.StaffSortBy(requested)
	if _, ok := allowed[sortBy]; ok {
		return sortBy.DBColumn()
	}
	return sortby.NewsSortByCreatedAt.DBColumn()
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

func (corporationService *CorporationService) GetCorporationStatuses() []corporationdto.GetEnumResponse {
	statuses := enum.GetAllCorporationStatuses()
	response := make([]corporationdto.GetEnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = corporationdto.GetEnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (corporationService *CorporationService) GetCorporationSortableColumns() []corporationdto.GetEnumResponse {
	columns := sortby.GetCorporationSortableColumns()
	response := make([]corporationdto.GetEnumResponse, len(columns))
	i := 0
	for col := range columns {
		response[i] = corporationdto.GetEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
}

func (corporationService *CorporationService) GetCorporationStaffSortableColumns() []corporationdto.GetEnumResponse {
	columns := sortby.GetStaffSortableColumns()
	response := make([]corporationdto.GetEnumResponse, len(columns))
	i := 0
	for col := range columns {
		response[i] = corporationdto.GetEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
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
		return nil, forbiddenError
	}
	return corporation, nil
}

func (corporationService *CorporationService) getCorporationByID(corporationID uint) (*entity.Corporation, error) {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return nil, err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return nil, notFoundError
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

func (corporationService *CorporationService) CheckApplicantAccess(corporationID, applicantID uint) error {
	staff, err := corporationService.corporationRepository.FindCorporationStaff(corporationService.db, applicantID, corporationID)
	if err != nil {
		return err
	}
	if staff == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return notFoundError
	}
	return nil
}

func (corporationService *CorporationService) ISCorporationApproved(corporationID uint) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if err != nil {
		return err
	}
	if corporation == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		return notFoundError
	}

	if corporation.Status != enum.CorpStatusApproved {
		return exception.NewUnapprovedCorporationForbiddenError()
	}
	return nil
}

func (corporationService *CorporationService) GetCorporationCredentials(corporationID uint) (corporationdto.CorporationCredentialResponse, error) {
	corporation, err := corporationService.getCorporationByID(corporationID)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
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
		Status:      corporation.Status.String(),
		ContactInfo: contactInfo,
		Addresses:   addresses,
	}, nil
}

func (corporationService *CorporationService) addSignatories(signatories []corporationdto.Signatory, corporationID uint) error {
	err := corporationService.db.WithTransaction(func(tx database.Database) error {
		for _, signatory := range signatories {
			signatoryModel, err := corporationService.corporationRepository.FindCorporationSignatoryByNationalID(corporationService.db, corporationID, signatory.NationalCardNumber, signatory.Position)
			if err != nil {
				return err
			}
			if signatoryModel != nil {
				continue
			}

			signatoryEntity := &entity.Signatory{
				CorporationID:      corporationID,
				Name:               signatory.Name,
				NationalCardNumber: signatory.NationalCardNumber,
				Position:           signatory.Position,
			}
			err = corporationService.corporationRepository.CreateSignatory(tx, signatoryEntity)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (corporationService *CorporationService) checkDuplicateName(conflictErrors *exception.ConflictErrors, name string, activeStatus []enum.CorporationStatus) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByName(corporationService.db, name, activeStatus)
	if err != nil {
		return err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.Name, corporationService.constants.Tag.AlreadyExist)
	}
	return nil
}

func (corporationService *CorporationService) checkDuplicateNationalID(conflictErrors *exception.ConflictErrors, nationalID string, activeStatus []enum.CorporationStatus) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByNationalID(corporationService.db, nationalID, activeStatus)
	if err != nil {
		return err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.NationalID, corporationService.constants.Tag.AlreadyExist)
	}
	return nil
}

func (corporationService *CorporationService) checkDuplicateRegistrationNumber(conflictErrors *exception.ConflictErrors, registrationNumber string, activeStatus []enum.CorporationStatus) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByRegistrationNumber(corporationService.db, registrationNumber, activeStatus)
	if err != nil {
		return err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.RegistrationNumber, corporationService.constants.Tag.AlreadyExist)
	}
	return nil
}

func (corporationService *CorporationService) checkDuplicateIBAN(conflictErrors *exception.ConflictErrors, iban string, activeStatus []enum.CorporationStatus) error {
	corporation, err := corporationService.corporationRepository.FindCorporationByIBAN(corporationService.db, iban, activeStatus)
	if err != nil {
		return err
	}
	if corporation != nil {
		conflictErrors.Add(corporationService.constants.Field.IBAN, corporationService.constants.Tag.AlreadyExist)
	}
	return nil
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) (corporationdto.CorporationCredentialResponse, error) {
	if err := corporationService.userService.IsUserActive(registerInfo.ApplicantID); err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	if err := corporationService.checkDuplicateName(&conflictErrors, registerInfo.Name, activeStatus); err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	if err := corporationService.checkDuplicateNationalID(&conflictErrors, registerInfo.NationalID, activeStatus); err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	if err := corporationService.checkDuplicateRegistrationNumber(&conflictErrors, registerInfo.RegistrationNumber, activeStatus); err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	if registerInfo.IBAN != "" {
		if err := corporationService.checkDuplicateIBAN(&conflictErrors, registerInfo.IBAN, activeStatus); err != nil {
			return corporationdto.CorporationCredentialResponse{}, err
		}
	}

	if len(conflictErrors.Errors) > 0 {
		return corporationdto.CorporationCredentialResponse{}, conflictErrors
	}

	role, err := corporationService.corporationRepository.FindRoleByName(corporationService.db, enum.CorporationOwner.String())
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	if role == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Role}
		return corporationdto.CorporationCredentialResponse{}, notFoundError
	}

	corporation := &entity.Corporation{
		Name:               registerInfo.Name,
		RegistrationNumber: registerInfo.RegistrationNumber,
		NationalID:         registerInfo.NationalID,
		IBAN:               registerInfo.IBAN,
		Status:             enum.CorpStatusAwaitingApproval,
	}

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		err := corporationService.corporationRepository.CreateCorporation(tx, corporation)
		if err != nil {
			return err
		}

		staff := &entity.CorporationStaff{
			UserID:        registerInfo.ApplicantID,
			CorporationID: corporation.ID,
			Status:        enum.StaffStatusActive,
			Roles:         []entity.Role{*role},
		}
		err = corporationService.corporationRepository.CreateCorporationStaff(tx, staff)
		if err != nil {
			return err
		}

		if err := corporationService.addSignatories(registerInfo.Signatories, corporation.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}

	return corporationdto.CorporationCredentialResponse{ID: corporation.ID, Name: corporation.Name}, nil
}

func (corporationService *CorporationService) replaceSignatories(corporationID uint, signatories []corporationdto.Signatory) error {
	err := corporationService.db.WithTransaction(func(tx database.Database) error {
		if err := corporationService.corporationRepository.DeleteCorporationSignatories(tx, corporationID); err != nil {
			return err
		}

		if err := corporationService.addSignatories(signatories, corporationID); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (corporationService *CorporationService) updateCorporationStatus(corporation *entity.Corporation) {
	if corporation.Status == enum.CorpStatusApproved {
		corporation.Status = enum.CorpStatusSuspend
	} else if corporation.Status == enum.CorpStatusRejected {
		corporation.Status = enum.CorpStatusAwaitingApproval
	}
}

func (corporationService *CorporationService) UpdateRegistrationInfoProfile(updateRegisterInfo corporationdto.UpdateRegisterRequest) error {
	corporation, err := corporationService.getCorporationByID(updateRegisterInfo.CorporationID)
	if err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(updateRegisterInfo.ApplicantID); err != nil {
		return err
	}

	if err := corporationService.CheckApplicantAccess(updateRegisterInfo.CorporationID, updateRegisterInfo.ApplicantID); err != nil {
		return err
	}

	if err := corporationService.checkCorporationConflicts(corporation, updateRegisterInfo.Name, updateRegisterInfo.NationalID, updateRegisterInfo.RegistrationNumber, updateRegisterInfo.IBAN); err != nil {
		return err
	}

	corporationService.updateCorporationStatus(corporation)

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		err = corporationService.corporationRepository.UpdateCorporation(tx, corporation)
		if err != nil {
			return err
		}

		err = corporationService.replaceSignatories(updateRegisterInfo.CorporationID, updateRegisterInfo.Signatories)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(updateRegisterInfo.CorporationID, enum.CorpStatusAwaitingApproval)
	if err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(updateRegisterInfo.ApplicantID); err != nil {
		return err
	}

	err = corporationService.CheckApplicantAccess(updateRegisterInfo.CorporationID, updateRegisterInfo.ApplicantID)
	if err != nil {
		return err
	}

	if err = corporationService.checkCorporationConflicts(corporation, updateRegisterInfo.Name, updateRegisterInfo.NationalID, updateRegisterInfo.RegistrationNumber, updateRegisterInfo.IBAN); err != nil {
		return err
	}

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		err = corporationService.corporationRepository.UpdateCorporation(tx, corporation)
		if err != nil {
			return err
		}

		err = corporationService.replaceSignatories(updateRegisterInfo.CorporationID, updateRegisterInfo.Signatories)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) checkCorporationConflicts(corporation *entity.Corporation, name, nationalID, registrationNumber, iban *string) error {
	activeStatus := []enum.CorporationStatus{enum.CorpStatusApproved, enum.CorpStatusAwaitingApproval}
	var conflictErrors exception.ConflictErrors
	if name != nil {
		if err := corporationService.checkDuplicateName(&conflictErrors, *name, activeStatus); err != nil {
			return err
		}
		corporation.Name = *name
	}

	if nationalID != nil {
		if err := corporationService.checkDuplicateNationalID(&conflictErrors, *nationalID, activeStatus); err != nil {
			return err
		}
		corporation.NationalID = *nationalID
	}

	if registrationNumber != nil {
		if err := corporationService.checkDuplicateRegistrationNumber(&conflictErrors, *registrationNumber, activeStatus); err != nil {
			return err
		}
		corporation.RegistrationNumber = *registrationNumber
	}

	if iban != nil {
		if err := corporationService.checkDuplicateIBAN(&conflictErrors, *iban, activeStatus); err != nil {
			return err
		}
		corporation.IBAN = *iban
	}

	if len(conflictErrors.Errors) > 0 {
		return conflictErrors
	}
	return nil
}

func (corporationService *CorporationService) AddCertificateFilesFromProfile(requestInfo corporationdto.AddCertificatesRequest) error {
	corporation, err := corporationService.getCorporationByID(requestInfo.CorporationID)
	if err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(requestInfo.ApplicantID); err != nil {
		return err
	}

	if err = corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.ApplicantID); err != nil {
		return err
	}

	prevVatTaxPayerPath := ""

	if requestInfo.VATTaxpayerCertificate != nil {
		taxPayerPath := corporationService.constants.S3BucketPath.GetVATTaxpayerCertificatePath(corporation.ID, requestInfo.VATTaxpayerCertificate.Filename)
		corporationService.s3Storage.UploadObject(enum.VATTaxpayerCertificate, taxPayerPath, requestInfo.VATTaxpayerCertificate)
		prevVatTaxPayerPath = corporation.VATTaxpayerCertificate
		corporation.VATTaxpayerCertificate = taxPayerPath
	}

	prevOfficialNewspaperPath := ""
	if requestInfo.OfficialNewspaperAD != nil {
		newspaperADPath := corporationService.constants.S3BucketPath.GetOfficialNewspaperADPath(corporation.ID, requestInfo.OfficialNewspaperAD.Filename)
		corporationService.s3Storage.UploadObject(enum.OfficialNewspaperAD, newspaperADPath, requestInfo.OfficialNewspaperAD)
		prevOfficialNewspaperPath = corporation.OfficialNewspaperAD
		corporation.OfficialNewspaperAD = newspaperADPath
	}

	corporationService.updateCorporationStatus(corporation)

	err = corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		return err
	}

	if prevVatTaxPayerPath != "" {
		corporationService.s3Storage.DeleteObject(enum.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate)
	}

	if prevOfficialNewspaperPath != "" {
		corporationService.s3Storage.DeleteObject(enum.OfficialNewspaperAD, corporation.OfficialNewspaperAD)
	}

	return nil
}

func (corporationService *CorporationService) AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(requestInfo.CorporationID, enum.CorpStatusAwaitingApproval)
	if err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(requestInfo.ApplicantID); err != nil {
		return err
	}

	if err = corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.ApplicantID); err != nil {
		return err
	}

	prevVatTaxPayerPath := ""

	if requestInfo.VATTaxpayerCertificate != nil {
		taxPayerPath := corporationService.constants.S3BucketPath.GetVATTaxpayerCertificatePath(corporation.ID, requestInfo.VATTaxpayerCertificate.Filename)
		corporationService.s3Storage.UploadObject(enum.VATTaxpayerCertificate, taxPayerPath, requestInfo.VATTaxpayerCertificate)
		prevVatTaxPayerPath = corporation.VATTaxpayerCertificate
		corporation.VATTaxpayerCertificate = taxPayerPath
	}

	prevOfficialNewspaperPath := ""
	if requestInfo.OfficialNewspaperAD != nil {
		newspaperADPath := corporationService.constants.S3BucketPath.GetOfficialNewspaperADPath(corporation.ID, requestInfo.OfficialNewspaperAD.Filename)
		corporationService.s3Storage.UploadObject(enum.OfficialNewspaperAD, newspaperADPath, requestInfo.OfficialNewspaperAD)
		prevOfficialNewspaperPath = corporation.OfficialNewspaperAD
		corporation.OfficialNewspaperAD = newspaperADPath
	}

	err = corporationService.corporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		return err
	}

	if prevVatTaxPayerPath != "" {
		corporationService.s3Storage.DeleteObject(enum.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate)
	}

	if prevOfficialNewspaperPath != "" {
		corporationService.s3Storage.DeleteObject(enum.OfficialNewspaperAD, corporation.OfficialNewspaperAD)
	}

	return nil
}

func (corporationService *CorporationService) AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) error {
	_, err := corporationService.getCorporationByIDAndStatus(contactInfo.CorporationID, contactInfo.CorporationStatus)
	if err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(contactInfo.ApplicantID); err != nil {
		return err
	}

	if err = corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID); err != nil {
		return err
	}

	for _, contact := range contactInfo.ContactInformation {
		contactType, err := corporationService.corporationRepository.FindContactInformationTypeByID(corporationService.db, contact.ContactTypeID)
		if err != nil {
			return err
		}
		if contactType == nil {
			continue
		}
		contactInfoType, err := corporationService.corporationRepository.FindContactInformationTypeValue(corporationService.db, contact.ContactTypeID, contact.ContactValue)
		if err != nil {
			return err
		}
		if contactInfoType != nil {
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

	if err := corporationService.userService.IsUserActive(contactInfo.ApplicantID); err != nil {
		return err
	}

	err = corporationService.CheckApplicantAccess(contactInfo.CorporationID, contactInfo.ApplicantID)
	if err != nil {
		return err
	}

	contact, err := corporationService.corporationRepository.FindContactInformationByID(corporationService.db, contactInfo.ContactID)
	if err != nil {
		return err
	}
	if contact == nil {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.ContactInformation}
		return notFoundError
	}

	err = corporationService.corporationRepository.DeleteContactInfo(corporationService.db, contact)
	if err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) getPrivateCorporationDetails(corporation *entity.Corporation) (corporationdto.CorporationPrivateInfoResponse, error) {
	vatTaxPayer := ""
	var err error
	if corporation.VATTaxpayerCertificate != "" {
		vatTaxPayer, err = corporationService.s3Storage.GetPresignedURL(enum.VATTaxpayerCertificate, corporation.VATTaxpayerCertificate, 8*time.Hour)
		if err != nil {
			return corporationdto.CorporationPrivateInfoResponse{}, err
		}
	}

	officialNewspaperAD := ""
	if corporation.OfficialNewspaperAD != "" {
		officialNewspaperAD, err = corporationService.s3Storage.GetPresignedURL(enum.OfficialNewspaperAD, corporation.OfficialNewspaperAD, 8*time.Hour)
		if err != nil {
			return corporationdto.CorporationPrivateInfoResponse{}, err
		}
	}

	logo := ""
	if corporation.Logo != "" {
		logo, err = corporationService.s3Storage.GetPresignedURL(enum.LogoPic, corporation.Logo, 8*time.Hour)
		if err != nil {
			return corporationdto.CorporationPrivateInfoResponse{}, err
		}
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
	corporation, err := corporationService.getCorporationByID(requestInfo.CorporationID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	if err = corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.UserID); err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	details, err := corporationService.getPrivateCorporationDetails(corporation)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	return details, nil
}

func (corporationService *CorporationService) GetCorporationPublicDetails(requestInfo corporationdto.CorporationDetailsRequest) (corporationdto.CorporationCredentialResponse, error) {
	if err := corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.UserID); err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	credentials, err := corporationService.GetCorporationCredentials(requestInfo.CorporationID)
	if err != nil {
		return corporationdto.CorporationCredentialResponse{}, err
	}
	return credentials, nil
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

	if err := corporationService.userService.IsUserActive(addressInfo.ApplicantID); err != nil {
		return err
	}

	if err = corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.ApplicantID); err != nil {
		return err
	}

	for _, address := range addressInfo.Addresses {
		corporationService.addressService.CreateAddress(address)
	}
	return nil
}

func (corporationService *CorporationService) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) error {
	if _, err := corporationService.getCorporationByIDAndStatus(addressInfo.CorporationID, addressInfo.CorporationStatus); err != nil {
		return err
	}

	if err := corporationService.userService.IsUserActive(addressInfo.UserID); err != nil {
		return err
	}

	if err := corporationService.CheckApplicantAccess(addressInfo.CorporationID, addressInfo.UserID); err != nil {
		return err
	}

	if err := corporationService.addressService.DeleteAddress(addressInfo.AddressID); err != nil {
		return err
	}

	return nil
}

func (corporationService *CorporationService) ChangeLogo(changeLogoRequest corporationdto.ChangeLogoRequest) error {
	corporation, err := corporationService.getCorporationByIDAndStatus(changeLogoRequest.CorporationID, enum.CorpStatusApproved)
	if err != nil {
		return err
	}

	if err = corporationService.CheckApplicantAccess(changeLogoRequest.CorporationID, changeLogoRequest.ApplicantID); err != nil {
		return err
	}

	prevLogoPath := ""
	if changeLogoRequest.Logo != nil {
		newLogoPath := corporationService.constants.S3BucketPath.GetCorporationLogoPath(changeLogoRequest.CorporationID, changeLogoRequest.Logo.Filename)
		corporationService.s3Storage.UploadObject(enum.LogoPic, newLogoPath, changeLogoRequest.Logo)
		corporation.Logo = newLogoPath
		prevLogoPath = corporation.Logo
	}

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		err = corporationService.corporationRepository.UpdateCorporation(tx, corporation)
		if err != nil {
			return err
		}

		if prevLogoPath != "" {
			if err := corporationService.s3Storage.DeleteObject(enum.LogoPic, corporation.Logo); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) GetUserCorporations(userID uint) ([]corporationdto.CorporationCredentialResponse, error) {
	corporations, err := corporationService.corporationRepository.FindUserActiveCorporations(corporationService.db, userID)
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
	corporations, err := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses, nil)
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

func (corporationService *CorporationService) getCorporationsByQuery(allowedStatuses []enum.CorporationStatus, query string, options *postgres.QueryOptions) ([]*entity.Corporation, int64, error) {
	if query == "" {
		corporations, err := corporationService.corporationRepository.FindCorporationsByStatus(corporationService.db, allowedStatuses, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := corporationService.corporationRepository.CountCorporationsByStatus(corporationService.db, allowedStatuses)
		if err != nil {
			return nil, 0, err
		}
		return corporations, count, nil
	}
	corporations, err := corporationService.corporationRepository.FindCorporationsByQuery(corporationService.db, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := corporationService.corporationRepository.CountCorporationsByQuery(corporationService.db, query)
	if err != nil {
		return nil, 0, err
	}
	return corporations, count, nil
}

func (corporationService *CorporationService) GetCorporationsByAdmin(listInfo corporationdto.GetCorporationsByAdminRequest) ([]corporationdto.CorporationCredentialResponse, int64, error) {
	allowedStatuses := corporationService.mapStatusIDToAllowedStatuses(listInfo.Status)

	options := postgres.NewQueryOptions().
		WithPagination(listInfo.Limit, listInfo.Offset).
		WithSorting(corporationService.getSortByColumn(listInfo.SortBy), listInfo.Asc)

	corporations, count, err := corporationService.getCorporationsByQuery(allowedStatuses, listInfo.Query, options)
	if err != nil {
		return nil, 0, err
	}

	response := make([]corporationdto.CorporationCredentialResponse, len(corporations))
	for i, corporation := range corporations {
		credentials, err := corporationService.GetCorporationCredentials(corporation.ID)
		if err != nil {
			return nil, 0, err
		}
		response[i] = credentials
	}

	return response, count, nil
}

func (corporationService *CorporationService) GetCorporationByAdmin(corporationID uint) (corporationdto.CorporationPrivateInfoResponse, error) {
	corporation, err := corporationService.getCorporationByID(corporationID)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}

	details, err := corporationService.getPrivateCorporationDetails(corporation)
	if err != nil {
		return corporationdto.CorporationPrivateInfoResponse{}, err
	}
	return details, nil
}

func (corporationService *CorporationService) GetReviewActions() []corporationdto.GetEnumResponse {
	actions := enum.GetAllReviewActions()
	response := make([]corporationdto.GetEnumResponse, len(actions))
	for i, action := range actions {
		response[i] = corporationdto.GetEnumResponse{
			ID:   uint(action),
			Name: action.String(),
		}
	}
	return response
}

func (corporationService *CorporationService) GetCorporationReviewsByAdmin(corporationID uint) ([]corporationdto.GetAdminCorporationReview, error) {
	reviews, err := corporationService.corporationRepository.FindCorporationReviews(corporationService.db, corporationID, nil)
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
	corporation, err := corporationService.getCorporationByID(request.CorporationID)
	if err != nil {
		return err
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

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		if err := corporationService.corporationRepository.CreateReview(tx, review); err != nil {
			return err
		}

		corporation.Status = enum.CorpStatusApproved
		if err := corporationService.corporationRepository.UpdateCorporation(tx, corporation); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) RejectCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error {
	corporation, err := corporationService.getCorporationByID(request.CorporationID)
	if err != nil {
		return err
	}

	if enum.ReviewAction(request.ActionID) == enum.ReviewActionApproved {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: corporationService.constants.Field.CorporationReview,
		}
		return forbiddenError
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

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		if err := corporationService.corporationRepository.CreateReview(tx, review); err != nil {
			return err
		}

		var corpStatus enum.CorporationStatus
		if enum.ReviewAction(request.ActionID) == enum.ReviewActionSuspended {
			corpStatus = enum.CorpStatusSuspend
		} else {
			corpStatus = enum.CorpStatusRejected
		}

		corporation.Status = corpStatus
		if err := corporationService.corporationRepository.UpdateCorporation(tx, corporation); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) GetStaffStatuses() []corporationdto.GetEnumResponse {
	statuses := enum.GetAllStaffStatuses()
	response := make([]corporationdto.GetEnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = corporationdto.GetEnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (corporationService *CorporationService) createStaff(UserID, CorporationID uint, roleIDs []uint) error {
	staff := &entity.CorporationStaff{
		UserID:        UserID,
		CorporationID: CorporationID,
		Status:        enum.StaffStatusActive,
	}

	updateStaffRolesRequest := rbacdto.UpdateStaffRolesRequest{
		Staff:   staff,
		RoleIDs: roleIDs,
	}
	err := corporationService.db.WithTransaction(func(tx database.Database) error {
		if err := corporationService.corporationRepository.CreateStaff(corporationService.db, staff); err != nil {
			return err
		}

		if err := corporationService.rbacService.UpdateStaffRoles(updateStaffRolesRequest); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (corporationService *CorporationService) updateStaffRoles(staff *entity.CorporationStaff, roleIDs []uint) error {
	updateStaffRolesRequest := rbacdto.UpdateStaffRolesRequest{
		Staff:   staff,
		RoleIDs: roleIDs,
	}

	if err := corporationService.rbacService.UpdateStaffRoles(updateStaffRolesRequest); err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) AddStaff(request corporationdto.AddStaffRequest) error {
	user, err := corporationService.userService.FindActiveUserByPhone(request.StaffPhone)
	if err != nil {
		return err
	}

	allowedStatuses := []enum.StaffStatus{enum.StaffStatusActive}
	staff, err := corporationService.corporationRepository.FindStaffByUserIDAndStatus(corporationService.db, user.ID, allowedStatuses)
	if err != nil {
		return err
	}
	if staff == nil {
		if err := corporationService.createStaff(user.ID, request.CorporationID, request.RoleIDs); err != nil {
			return err
		}
		return nil
	}

	if err := corporationService.updateStaffRoles(staff, request.RoleIDs); err != nil {
		return err
	}
	return nil
}

func (corporationService *CorporationService) getCorporationStaffByID(corporationID, staffID uint) (*entity.CorporationStaff, error) {
	staff, err := corporationService.corporationRepository.FindCorporationStaffByID(corporationService.db, corporationID, staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, exception.NotFoundError{Item: corporationService.constants.Field.CorporationStaff}
	}
	return staff, nil
}

func (corporationService *CorporationService) mapToStaffStatus(requested uint) []enum.StaffStatus {
	if requested == uint(enum.StaffStatusAll) {
		return enum.GetAllStaffStatuses()
	}
	for _, status := range enum.GetAllStaffStatuses() {
		if requested == uint(status) && status != enum.StaffStatusAll {
			return []enum.StaffStatus{status}
		}
	}
	return enum.GetAllStaffStatuses()
}

func (corporationService *CorporationService) EditStaff(request corporationdto.EditStaffRequest) error {
	staff, err := corporationService.getCorporationStaffByID(request.CorporationID, request.StaffID)
	if err != nil {
		return err
	}

	if staff.Status != enum.StaffStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(corporationService.constants.Field.CorporationStaff, corporationService.constants.Tag.NotActive)
		return conflictErrors
	}

	if request.Status != nil {
		if *request.Status == uint(enum.StaffStatusInactive) {
			staff.Status = enum.StaffStatusInactive
		} else if *request.Status == uint(enum.StaffStatusActive) {
			staff.Status = enum.StaffStatusActive
		}
	}

	err = corporationService.db.WithTransaction(func(tx database.Database) error {
		if err := corporationService.updateStaffRoles(staff, request.RoleIDs); err != nil {
			return err
		}

		if err := corporationService.corporationRepository.UpdateStaff(corporationService.db, staff); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (corporationService *CorporationService) getStaffsByQuery(corporationID uint, allowedStatus []enum.StaffStatus, query string, options *postgres.QueryOptions) ([]*entity.CorporationStaff, int64, error) {
	if query == "" {
		staffs, err := corporationService.corporationRepository.FindCorporationStaffs(corporationService.db, corporationID, allowedStatus, options)
		if err != nil {
			log.Println("hi am here")
			return nil, 0, err
		}
		count, err := corporationService.corporationRepository.CountCorporationStaffs(corporationService.db, corporationID, allowedStatus)
		if err != nil {
			return nil, 0, err
		}
		return staffs, count, nil
	}

	staffs, err := corporationService.corporationRepository.FindCorporationStaffByQuery(corporationService.db, corporationID, allowedStatus, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := corporationService.corporationRepository.CountCorporationStaffByQuery(corporationService.db, corporationID, allowedStatus, query)
	if err != nil {
		return nil, 0, err
	}
	return staffs, count, nil
}

func (corporationService *CorporationService) GetStaffList(request corporationdto.GetStaffList) ([]corporationdto.StaffDetailsResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(corporationService.getStaffSortByColumn(request.SortBy), request.Asc)

	allowedStatus := corporationService.mapToStaffStatus(request.Status)
	staffs, count, err := corporationService.getStaffsByQuery(request.CorporationID, allowedStatus, request.Query, options)
	if err != nil {
		return nil, 0, err
	}
	response := make([]corporationdto.StaffDetailsResponse, len(staffs))

	for i, staff := range staffs {
		user, err := corporationService.userService.GetUserCredential(staff.UserID)
		if err != nil {
			return nil, 0, err
		}

		roles, err := corporationService.rbacService.GetStaffRoles(staff)
		if err != nil {
			return nil, 0, err
		}

		response[i] = corporationdto.StaffDetailsResponse{
			ID:     staff.ID,
			Staff:  user,
			Roles:  roles,
			Status: staff.Status.String(),
		}
	}
	return response, count, nil
}

func (corporationService *CorporationService) GetStaff(corporationID, staffID uint) (corporationdto.StaffDetailsResponse, error) {
	staff, err := corporationService.getCorporationStaffByID(corporationID, staffID)
	if err != nil {
		return corporationdto.StaffDetailsResponse{}, err
	}

	user, err := corporationService.userService.GetUserCredential(staff.UserID)
	if err != nil {
		return corporationdto.StaffDetailsResponse{}, err
	}

	roles, err := corporationService.rbacService.GetStaffRoles(staff)
	if err != nil {
		return corporationdto.StaffDetailsResponse{}, err
	}
	return corporationdto.StaffDetailsResponse{
		ID:     staff.ID,
		Staff:  user,
		Roles:  roles,
		Status: staff.Status.String(),
	}, nil
}
