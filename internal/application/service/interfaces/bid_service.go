package service

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	GetBidStatuses() []biddto.GetBidStatusesResponse
	AcceptBid(request biddto.GetCustomerBidRequest)
	CancelBid(bidInfo biddto.GetBidRequest)
	GetCorporationBid(request biddto.GetBidRequest) biddto.CorporationBidResponse
	GetCorporationBids(request biddto.GetCorporationBidsRequest) []biddto.CorporationBidResponse
	GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) biddto.AnonymousBidResponse
	GetRequestAnonymousBids(requestInfo biddto.GetListRequestBidsRequest) []biddto.AnonymousBidResponse
	GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) []biddto.AdminBidResponse
	RejectBid(request biddto.GetCustomerBidRequest)
	SetBid(bidInfo biddto.SetBidRequest)
	UpdateBid(request biddto.UpdateBidRequest)
}
