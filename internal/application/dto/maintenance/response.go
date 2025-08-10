package maintenancedto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type MaintenanceEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CustomerMaintenanceRequestListResponse struct {
	ID                   uint                                         `json:"id"`
	CreatedAt            time.Time                                    `json:"createdAt"`
	Panel                installationdto.CustomerPanelResponse        `json:"panel"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Subject              string                                       `json:"subject"`
	Description          string                                       `json:"description"`
	UrgencyLevel         string                                       `json:"urgencyLevel"`
	Status               string                                       `json:"status"`
	IsGuaranteeRequested bool                                         `json:"isGuaranteeRequested"`
}

type CustomerMaintenanceRequestResponse struct {
	ID                   uint                                         `json:"id"`
	CreatedAt            time.Time                                    `json:"createdAt"`
	Panel                installationdto.CustomerPanelResponse        `json:"panel"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Subject              string                                       `json:"subject"`
	Description          string                                       `json:"description"`
	UrgencyLevel         string                                       `json:"urgencyLevel"`
	Status               string                                       `json:"status"`
	IsGuaranteeRequested bool                                         `json:"isGuaranteeRequested"`
	Record               CustomerMaintenanceRecordResponse            `json:"record"`
}

type CustomerMaintenanceRecordResponse struct {
	ID                 uint                                            `json:"id"`
	CreatedAt          time.Time                                       `json:"createdAt"`
	Title              string                                          `json:"title"`
	Details            string                                          `json:"details"`
	Date               time.Time                                       `json:"date"`
	IsUserApproved     bool                                            `json:"isApproved"`
	GuaranteeViolation guaranteedto.CustomerGuaranteeViolationResponse `json:"violation"`
}

type CorporationMaintenanceListResponse struct {
	ID                   uint                                     `json:"id"`
	CreatedAt            time.Time                                `json:"createdAt"`
	Panel                installationdto.CorporationPanelResponse `json:"panel"`
	Subject              string                                   `json:"subject"`
	Description          string                                   `json:"description"`
	UrgencyLevel         string                                   `json:"urgencyLevel"`
	Status               string                                   `json:"status"`
	IsGuaranteeRequested bool                                     `json:"isGuaranteeRequested"`
}

type CorporationMaintenanceResponse struct {
	ID                   uint                                     `json:"id"`
	CreatedAt            time.Time                                `json:"createdAt"`
	Panel                installationdto.CorporationPanelResponse `json:"panel"`
	Subject              string                                   `json:"subject"`
	Description          string                                   `json:"description"`
	UrgencyLevel         string                                   `json:"urgencyLevel"`
	Status               string                                   `json:"status"`
	IsGuaranteeRequested bool                                     `json:"isGuaranteeRequested"`
	Record               CorporationMaintenanceRecordResponse     `json:"record"`
}

type CorporationMaintenanceRecordResponse struct {
	ID                 uint                                               `json:"id"`
	CreatedAt          time.Time                                          `json:"createdAt"`
	Operator           userdto.CredentialResponse                         `json:"operator"`
	Title              string                                             `json:"title"`
	Details            string                                             `json:"details"`
	IsUserApproved     bool                                               `json:"isApproved"`
	GuaranteeViolation guaranteedto.CorporationGuaranteeViolationResponse `json:"violation"`
}

type AdminMaintenanceRequestResponse struct {
	ID                   uint                                         `json:"id"`
	CreatedAt            time.Time                                    `json:"createdAt"`
	Panel                installationdto.AdminPanelResponse           `json:"panel"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Subject              string                                       `json:"subject"`
	Description          string                                       `json:"description"`
	UrgencyLevel         string                                       `json:"urgencyLevel"`
	Status               string                                       `json:"status"`
	IsGuaranteeRequested bool                                         `json:"isGuaranteeRequested"`
	Record               CorporationMaintenanceRecordResponse         `json:"record"`
}
