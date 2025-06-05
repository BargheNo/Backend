package biddto

import (
	"time"

	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

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

type GetBidStatusesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
