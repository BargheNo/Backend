package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationLoginResponse
	ChangePassword(changePasswordRequest corporationdto.ChangePasswordRequest)
	GetCorporationInfo(idRequest corporationdto.IDRequest) corporationdto.CorporationInfoResponse
	UpdateContactInfo(corporationID uint, contactInfo corporationdto.ContactInfoRequest)
	AddAddress(address corporationdto.AddressRequest)
	EditAddress(addressID uint, address corporationdto.AddressRequest)
	DeleteAddress(corporationID uint, addressID uint)
	GetCorporationByID(corporationID uint) (*entity.Corporation, bool)
}
