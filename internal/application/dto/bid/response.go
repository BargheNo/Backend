package biddto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type InstallationRequestDetails struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CustomerName string                     `json:"customerName"`
	Address      addressdto.AddressResponse `json:"address"`
	PowerRequest uint                       `json:"powerRequest"`
}

type AnonymousBidResponse struct {
	ID               uint                            `json:"id"`
	Description      string                          `json:"description"`
	Cost             uint                            `json:"cost"`
	InstallationTime time.Time                       `json:"installationTime"`
	Status           string                          `json:"status"`
	PaymentTerms     paymentdto.PaymentTermsResponse `json:"paymentTerms"`
	Guarantee        guaranteedto.GuaranteeResponse  `json:"guarantee"`
}

type CorporationBidResponse struct {
	ID                  uint                                      `json:"id"`
	Bidder              userdto.CredentialResponse                `json:"bidder"`
	InstallationRequest installationdto.AnonymousRequestsResponse `json:"request"`
	Description         string                                    `json:"description"`
	Cost                uint                                      `json:"cost"`
	InstallationTime    time.Time                                 `json:"installationTime"`
	Status              string                                    `json:"status"`
	PaymentTerms        paymentdto.PaymentTermsResponse           `json:"paymentTerms"`
	Guarantee           guaranteedto.GuaranteeResponse            `json:"guarantee"`
}
