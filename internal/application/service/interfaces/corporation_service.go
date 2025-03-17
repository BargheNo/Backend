package service

import corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse
	GetInstallationRequests(id uint) []corporationdto.InstallationRequestResponse
}
