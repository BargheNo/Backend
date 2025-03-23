package installationdto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type ListOwnerRequestsResponse struct {
	ID          uint                       `json:"id"`
	Name        string                     `json:"name"`
	Status      string                     `json:"status"`
	CreatedTime time.Time                  `json:"createdTime"`
	Address     addressdto.AddressResponse `json:"address"`
}

type InstallationRequestDetails struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CustomerName string                     `json:"customerName"`
	Address      addressdto.AddressResponse `json:"address"`
	PowerRequest uint                       `json:"powerRequest"`
}
