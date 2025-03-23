package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type CorporationService interface {
	GetCorporationByID(corporationID uint) (*entity.Corporation, bool)
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationLoginResponse
	ChangePassword(changePasswordRequest corporationdto.ChangePasswordRequest)
	UpdateContactInfo(contactInfo corporationdto.ContactInfoRequest)
	GetCorporationInfo(idRequest corporationdto.IDRequest) corporationdto.CorporationInfoResponse
	// AddAddress(address corporationdto.AddressRequest)
	// EditAddress(addressID uint, address corporationdto.AddressRequest)
	// DeleteAddress(corporationID uint, addressID uint)
}
