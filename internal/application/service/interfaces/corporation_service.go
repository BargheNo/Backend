package service

import corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse
	UpdateContactInfo(corporationID uint, contactInfo corporationdto.ContactInfoRequest)
	AddAddress(corporationID uint, address corporationdto.AddressRequest)
	EditAddress(corporationID uint, addressID uint, address corporationdto.AddressRequest)
	DeleteAddress(corporationID uint, addressID uint)
}
