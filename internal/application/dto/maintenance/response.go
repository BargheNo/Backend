package maintenancedto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type CustomerMaintenanceRequestResponse struct {
	ID           uint                                         `json:"id"`
	Panel        installationdto.CustomerPanelResponse        `json:"panel"`
	Corporation  corporationdto.CorporationCredentialResponse `json:"corporation"`
	OwnerID      uint                                         `json:"ownerID"`
	Subject      string                                       `json:"subject"`
	Description  string                                       `json:"description"`
	UrgencyLevel string                                       `json:"urgencyLevel"`
	Status       string                                       `json:"status"`
	CreatedAt    time.Time                                    `json:"createdAt"`
}

type CorporationMaintenanceResponse struct {
	ID           uint                                     `json:"id"`
	Panel        installationdto.CorporationPanelResponse `json:"panel"`
	Subject      string                                   `json:"subject"`
	Description  string                                   `json:"description"`
	UrgencyLevel string                                   `json:"urgencyLevel"`
	Status       string                                   `json:"status"`
	CreatedAt    time.Time                                `json:"createdAt"`
	OwnerPhone   string                                   `json:"ownerPhone"`
}

type MaintenanceRequestResponse struct {
	ID           uint                                         `json:"id"`
	Panel        installationdto.PanelResponse                `json:"panel"`
	Corporation  corporationdto.CorporationCredentialResponse `json:"corporation"`
	Customer     userdto.CredentialResponse                   `json:"customer"`
	Subject      string                                       `json:"subject"`
	Description  string                                       `json:"description"`
	UrgencyLevel string                                       `json:"urgencyLevel"`
	Status       string                                       `json:"status"`
	CreatedAt    time.Time                                    `json:"createdAt"`
}

type CorporationMaintenanceRecordResponse struct {
	ID       uint                                     `json:"id"`
	Panel    installationdto.CorporationPanelResponse `json:"panel"`
	Operator userdto.CredentialResponse               `json:"operator"`
	Title    string                                   `json:"title"`
	Details  string                                   `json:"details"`
	Date     time.Time                                `json:"date"`
}

type CustomerMaintenanceRecordResponse struct {
	ID          uint                                         `json:"id"`
	Panel       installationdto.CustomerPanelResponse        `json:"panel"`
	Corporation corporationdto.CorporationCredentialResponse `json:"corporation"`
	Operator    userdto.CredentialResponse                   `json:"operator"`
	Title       string                                       `json:"title"`
	Details     string                                       `json:"details"`
	Date        time.Time                                    `json:"date"`
}

type MaintenanceRecordResponse struct {
	ID          uint                                         `json:"id"`
	Panel       installationdto.PanelResponse                `json:"panel"`
	Corporation corporationdto.CorporationCredentialResponse `json:"corporation"`
	Customer    userdto.CredentialResponse                   `json:"customer"`
	Operator    userdto.CredentialResponse                   `json:"operator"`
	Title       string                                       `json:"title"`
	Details     string                                       `json:"details"`
	Date        time.Time                                    `json:"date"`
}
