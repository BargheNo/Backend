package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse
	ChangePassword(changePasswordRequest corporationdto.ChangePasswordRequest)
	UpdateContactInfo(corporationID uint, contactInfo corporationdto.ContactInfoRequest)
	AddAddress(address corporationdto.AddressRequest)
	EditAddress(addressID uint, address corporationdto.AddressRequest)
	DeleteAddress(corporationID uint, addressID uint)
	GetCorporationByID(corporationID uint) (*entity.Corporation, bool)
}
