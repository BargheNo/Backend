package biddto

import "time"

type SetBidRequest struct {
	CorporationID         uint
	BidderID              uint
	InstallationRequestID uint
	Cost                  uint
	Description           string
	InstallationTime      time.Time
}

type CancelBidRequest struct {
	CorporationID         uint
	BidderID              uint
	BidID                 uint
	InstallationRequestID uint
}

type GetCorporationBidsRequest struct {
	CorporationID uint
	UserID        uint
	Offset        int
	Limit         int
}

type GetRequestBidsRequest struct {
	RequestID uint
	UserID    uint
}

type BidNotificationData struct {
	BidID                 uint
}
