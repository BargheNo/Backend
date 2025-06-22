package port

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
)

type CorporationService interface {
	GetCorporationStatuses() []corporationdto.GetStatusesResponse
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
	GetAvailableCorporations() ([]corporationdto.CorporationCredentialResponse, error)
	GetCorporationsByAdmin(listInfo corporationdto.GetCorporationsByAdminRequest) ([]corporationdto.CorporationCredentialResponse, error)
	GetCorporationByAdmin(corporationID uint) (corporationdto.CorporationPrivateInfoResponse, error)
	GetReviewActions() []corporationdto.GetStatusesResponse
	GetCorporationReviewsByAdmin(corporationID uint) ([]corporationdto.GetAdminCorporationReview, error)
	ApproveCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error
	RejectCorporationRegistration(request corporationdto.HandleCorporationActionRequest) error
}
