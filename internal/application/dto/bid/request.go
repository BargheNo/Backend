package biddto

import "time"

type SetBidRequest struct {
	InstallationRequestID uint
	CorporationID         uint
	MinCost               float64
	MaxCost               float64
	MinDeadline           time.Time
	MaxDeadline           time.Time
	Description           string
	InstallationTime      string
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
	SortBy        string
	Dir           string
}
