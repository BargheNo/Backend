package biddto

import "time"

// For corporation's guarantee template
type GuaranteeTemplateResponse struct {
	ID              uint      `json:"id"`
	GuaranteeType   string    `json:"guaranteeType"`
	DurationMonths  uint      `json:"durationMonths"`
	Description     string    `json:"description"`
	Terms           string    `json:"terms"`
	CoverageDetails string    `json:"coverageDetails"`
	IsActive        bool      `json:"isActive"`
}

// For selecting a guarantee when creating/updating a bid
type BidGuaranteeRequest struct {
	GuaranteeID uint      `json:"guaranteeId" validate:"required"`
	StartDate   time.Time `json:"startDate" validate:"required"`
}

// For viewing guarantee details in a bid
type BidGuaranteeResponse struct {
	ID            uint      `json:"id"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	IsActive      bool      `json:"isActive"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	Guarantee     GuaranteeResponse `json:"guarantee"`
}

// For viewing available guarantee templates
type GuaranteeResponse struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	GuaranteeType string              `json:"guaranteeType"`
	DurationMonths uint               `json:"durationMonths"`
	Description   string              `json:"description"`
	Terms         []GuaranteeTermResponse `json:"terms"`
}

type GuaranteeTermResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GuaranteeCoverageResponse struct {
	Item        string `json:"item"`
	Details     string `json:"details"`
	Limitations string `json:"limitations,omitempty"`
} 