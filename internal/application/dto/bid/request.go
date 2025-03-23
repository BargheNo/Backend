package biddto

import "time"

type SetBidRequest struct {
	InstallationRequestID uint
	CorporationID         uint
	Cost                  uint
	Description           string
	InstallationDate      time.Time
}

type CancelBidRequest struct {
	BidID                 uint
	InstallationRequestID uint
	CorporationID         uint
}

type GetBidsRequest struct {
	CorporationID uint
	Offset        int
	Limit         int
}
