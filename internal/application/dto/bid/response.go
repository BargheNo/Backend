package biddto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
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
	Area             uint                            `json:"area"`
	Power            uint                            `json:"power"`
	PaymentTerms     paymentdto.PaymentTermsResponse `json:"paymentTerms"`
	Guarantee        guaranteedto.GuaranteeResponse  `json:"guarantee"`
}

type CorporationBidResponse struct {
	ID                  uint                                      `json:"id"`
	Bidder              userdto.CredentialResponse                `json:"bidder"`
	InstallationRequest installationdto.AnonymousRequestsResponse `json:"request"`
	Description         string                                    `json:"description"`
	Cost                uint                                      `json:"cost"`
	Area                uint                                      `json:"area"`
	Power               uint                                      `json:"power"`
	InstallationTime    time.Time                                 `json:"installationTime"`
	Status              string                                    `json:"status"`
	PaymentTerms        paymentdto.PaymentTermsResponse           `json:"paymentTerms"`
	Guarantee           guaranteedto.GuaranteeResponse            `json:"guarantee"`
}

type AdminBidResponse struct {
	ID               uint                                         `json:"id"`
	Corporation      corporationdto.CorporationCredentialResponse `json:"corporation"`
	Bidder           userdto.CredentialResponse                   `json:"bidder"`
	Description      string                                       `json:"description"`
	Cost             uint                                         `json:"cost"`
	Area             uint                                         `json:"area"`
	Power            uint                                         `json:"power"`
	InstallationTime time.Time                                    `json:"installationTime"`
	Status           string                                       `json:"status"`
	PaymentTerms     paymentdto.PaymentTermsResponse              `json:"paymentTerms"`
	Guarantee        guaranteedto.GuaranteeResponse               `json:"guarantee"`
}

type GetBidEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
