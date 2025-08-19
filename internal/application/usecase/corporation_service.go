package usecase

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type CorporationService interface {
	GetCorporationSortableColumns() []corporationdto.GetEnumResponse
	GetCorporationStatuses() []corporationdto.GetEnumResponse
	DoesCorporationExist(corporationID uint) error
	ISCorporationApproved(corporationID uint) error
	GetCorporationCredentials(corporationID uint) (corporationdto.CorporationCredentialResponse, error)
	CheckApplicantAccess(corporationID, applicantID uint) error
	Register(registerInfo corporationdto.RegisterRequest) (corporationdto.CorporationCredentialResponse, error)
	UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest) error
	AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) error
	AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) error
	DeleteContactInfo(contactInfo corporationdto.DeleteContactInformationRequest) error
	AddAddress(addressInfo corporationdto.AddCorporationAddressRequest) error
	DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) error
	GetCorporationDetails(requestInfo corporationdto.CorporationDetailsRequest) (corporationdto.CorporationPrivateInfoResponse, error)
	GetContactTypes() ([]corporationdto.ContactTypeResponse, error)
	ChangeLogo(changeLogoRequest corporationdto.ChangeLogoRequest) error
	GetUserCorporations(userID uint) ([]corporationdto.CorporationCredentialResponse, error)
	UpdateRegistrationInfoProfile(updateRegisterInfo corporationdto.UpdateRegisterRequest) error
	AddCertificateFilesFromProfile(requestInfo corporationdto.AddCertificatesRequest) error
	GetAvailableCorporations() ([]corporationdto.CorporationCredentialResponse, error)
	GetCorporationsByAdmin(listInfo corporationdto.GetCorporationsByAdminRequest) ([]corporationdto.CorporationCredentialResponse, int64, error)
	GetCorporationPublicDetails(requestInfo corporationdto.CorporationDetailsRequest) (corporationdto.CorporationCredentialResponse, error)
	GetCorporationByAdmin(corporationID uint) (corporationdto.CorporationPrivateInfoResponse, error)
	GetReviewActions() []corporationdto.GetEnumResponse
	GetCorporationReviewsByAdmin(corporationID uint) ([]corporationdto.GetAdminCorporationReview, error)
	ApproveCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error
	RejectCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error
	GetStaffStatuses() []corporationdto.GetEnumResponse
	AddStaff(request corporationdto.AddStaffRequest) error
	EditStaff(request corporationdto.EditStaffRequest) error
	GetStaffList(request corporationdto.GetStaffList) ([]corporationdto.StaffDetailsResponse, int64, error)
	GetStaff(corporationID, staffID uint) (corporationdto.StaffDetailsResponse, error)
	GetCorporationRoles(request corporationdto.GetRolesListRequest) ([]userdto.RoleResponse, int64, error)
}
