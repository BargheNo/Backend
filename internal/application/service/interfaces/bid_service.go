package service

import "github.com/BargheNo/Backend/internal/application/dto/corporation"

type BidService interface {
	GetInstallationRequests(corporationID uint, page int, pageSize int) []corporationdto.InstallationRequestResponse
	SetBid(bidInfo corporationdto.SetBidRequest)
	CancelBid(bidInfo corporationdto.CancelBidRequest)
	GetBids(corporationID uint, page int, pageSize int) []corporationdto.BidsResponse
}
