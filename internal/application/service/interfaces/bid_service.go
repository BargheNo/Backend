package service

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	GetInstallationRequests(corporationID uint, page int, pageSize int, sortBy string, dir string) []biddto.InstallationRequestResponse
	SetBid(bidInfo biddto.SetBidRequest)
	CancelBid(bidInfo biddto.CancelBidRequest)
	GetBids(corporationID uint, page int, pageSize int, sortBy string, dir string) []biddto.BidsResponse
}
