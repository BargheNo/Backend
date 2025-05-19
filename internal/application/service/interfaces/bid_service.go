package service

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	SetBid(bidInfo biddto.SetBidRequest)
	CancelBid(bidInfo biddto.CancelBidRequest)
	GetBid(bidID uint) biddto.BidsResponse
	GetCorporationBids(bidsRequest biddto.GetCorporationBidsRequest) []biddto.BidsResponse
	GetRequestBids(requestInfo biddto.GetRequestBidsRequest) []biddto.BidsResponse
}
