package service

import corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest)
	Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse
	GetInstallationRequests(corporationID uint, page int, pageSize int) []corporationdto.InstallationRequestResponse
	SetBid(bidInfo corporationdto.SetBidRequest)
	CancelBid(bidInfo corporationdto.CancelBidRequest)
	GetBids(corporationID uint, page int, pageSize int) []corporationdto.BidsResponse
}
