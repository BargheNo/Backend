package service

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	SetBid(bidInfo biddto.SetBidRequest)
	CancelBid(bidInfo biddto.CancelBidRequest)
	GetBids(bidsRequest biddto.GetBidsRequest) []biddto.BidsResponse
}
