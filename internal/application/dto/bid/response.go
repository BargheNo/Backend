package biddto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
)

type InstallationRequestDetails struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CustomerName string                     `json:"customerName"`
	Address      addressdto.AddressResponse `json:"address"`
	PowerRequest uint                       `json:"powerRequest"`
}

type BidsResponse struct {
	ID                         uint                                   `json:"id"`
	InstallationRequestDetails installationdto.RequestDetailsResponse `json:"installationRequest"`
	Description                string                                 `json:"description"`
	Cost                       uint                                   `json:"cost"`
	InstallationDate           time.Time                              `json:"installationTime"`
	Status                     string                                 `json:"status"`
}
