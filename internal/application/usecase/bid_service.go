package usecase

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	GetBidStatuses() []biddto.GetBidStatusesResponse
	AcceptBid(request biddto.GetCustomerBidRequest) error
	CancelBid(bidInfo biddto.GetBidRequest) error
	GetCorporationBid(request biddto.GetBidRequest) (biddto.CorporationBidResponse, error)
	GetCorporationBids(request biddto.GetCorporationBidsRequest) ([]biddto.CorporationBidResponse, error)
	GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) (biddto.AnonymousBidResponse, error)
	GetRequestAnonymousBids(requestInfo biddto.GetListRequestBidsRequest) ([]biddto.AnonymousBidResponse, error)
	GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) ([]biddto.AdminBidResponse, error)
	RejectBid(request biddto.GetCustomerBidRequest) error
	SetBid(bidInfo biddto.SetBidRequest) error
	UpdateBid(request biddto.UpdateBidRequest) error
}
