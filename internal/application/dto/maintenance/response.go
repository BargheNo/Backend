package maintenancedto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type CustomerMaintenanceRequestResponse struct {
	ID           uint
	Panel        installationdto.CustomerPanelResponse
	Corporation  corporationdto.CorporationDetailsResponse
	OwnerID      uint
	Subject      string
	Description  string
	UrgencyLevel string
	Status       string
	CreatedAt    time.Time
}

type CorporationMaintenanceResponse struct {
	ID           uint
	Panel        installationdto.CorporationPanelResponse
	Subject      string
	Description  string
	UrgencyLevel string
	Status       string
	CreatedAt    time.Time
	OwnerPhone   string
}

type MaintenanceRequestResponse struct {
	ID           uint
	Panel        installationdto.PanleResponse
	Corporation  corporationdto.CorporationDetailsResponse
	Customer     userdto.CredentialResponse
	Subject      string
	Description  string
	UrgencyLevel string
	Status       string
	CreatedAt    time.Time
}

type CorporationMaintenanceRecordResponse struct {
	ID        uint
	RequestID uint
	Panel     installationdto.CorporationPanelResponse
	Operator  userdto.CredentialResponse
	Title     string
	Details   string
	Date      time.Time
}

type CustomerMaintenanceRecordResponse struct {
	ID          uint
	Panel       installationdto.CustomerPanelResponse
	Corporation corporationdto.CorporationDetailsResponse
	Operator    userdto.CredentialResponse
	Title       string
	Details     string
	Date        time.Time
}

type MaintenanceRecordResponse struct {
	ID          uint
	Panel       installationdto.PanleResponse
	Corporation corporationdto.CorporationDetailsResponse
	Customer    userdto.CredentialResponse
	Operator    userdto.CredentialResponse
	Title       string
	Details     string
	Date        time.Time
}
