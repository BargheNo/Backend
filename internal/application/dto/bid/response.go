package biddto

import "time"

type InstallationRequestResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"userId"`
	Area           float64   `json:"area"`
	PowerRequested float64   `json:"powerRequested"`
	MaxCost        float64   `json:"maxCost"`
	Deadline       time.Time `json:"deadline"`
	BuildingType   string    `json:"buildingType"`
	Address        string    `json:"address"`
}

type BidsResponse struct {
	ID                    uint      `json:"id"`
	InstallationRequestID uint      `json:"installationRequestId"`
	Description           string    `json:"description"`
	MinCost               float64   `json:"minCost"`
	MaxCost               float64   `json:"maxCost"`
	MinDeadline           time.Time `json:"minDeadline"`
	MaxDeadline           time.Time `json:"maxDeadline"`
	InstallationTime      string    `json:"installationTime"`
	Status                string    `json:"status"`
}
