package usecase

import (
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
)

type BidService interface {
	GetBidSortableColumns() []biddto.GetBidEnumResponse
	GetCorporationBidStatuses() []biddto.GetBidEnumResponse
	GetUserBidStatuses() []biddto.GetBidEnumResponse
	AcceptBid(request biddto.GetCustomerBidRequest) error
	CancelBid(bidInfo biddto.GetBidRequest) error
	GetCorporationBid(request biddto.GetBidRequest) (biddto.CorporationBidResponse, error)
	GetCorporationBids(request biddto.GetCorporationBidsRequest) ([]biddto.CorporationBidResponse, int64, error)
	GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) (biddto.AnonymousBidResponse, error)
	GetRequestAnonymousBids(requestInfo biddto.GetListRequestBidsRequest) ([]biddto.AnonymousBidResponse, int64, error)
	GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) ([]biddto.AdminBidResponse, int64, error)
	RejectBid(request biddto.GetCustomerBidRequest) error
	SetBid(bidInfo biddto.SetBidRequest) error
	UpdateBid(request biddto.UpdateBidRequest) error
	DeleteBidByAdmin(bidID uint) error
	UpdateBidByAdmin(request biddto.UpdateBidRequest) error
	GetBidsByAdmin(requestInfo biddto.GetListBidsRequestByAdmin) ([]biddto.AdminBidResponse, int64, error)
	GetBidByAdmin(bidID uint) (biddto.AdminBidResponse, error)
}
