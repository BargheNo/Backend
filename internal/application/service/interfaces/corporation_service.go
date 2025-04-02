package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type CorporationService interface {
	GetCorporationByID(corporationID uint) *entity.Corporation
	CheckApplicantAccess(corporationID, applicantID uint)
	Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationDetailsResponse
	AddContactInfo(contactInfo corporationdto.AddContactInformationRequest)
	AddAddress(addressInfo corporationdto.AddCorporationAddressRequest)
	DeleteAddress(addressInfo corporationdto.DeleteAddressRequest)
}
