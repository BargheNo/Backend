package biddto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type InstallationRequestDetails struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CustomerName string                     `json:"customerName"`
	Address      addressdto.AddressResponse `json:"address"`
	PowerRequest uint                       `json:"powerRequest"`
}

type BidsResponse struct {
	ID                         uint                       `json:"id"`
	InstallationRequestDetails InstallationRequestDetails `json:"installationRequestId"`
	Description                string                     `json:"description"`
	Cost                       uint                       `json:"cost"`
	InstallationDate           time.Time                  `json:"installationTime"`
	Status                     string                     `json:"status"`
}
