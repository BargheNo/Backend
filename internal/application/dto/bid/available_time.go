package biddto

import "time"

type AvailableTimeRequest struct {
	StartTime time.Time `json:"startTime" validate:"required"`
	EndTime   time.Time `json:"endTime" validate:"required,gtfield=StartTime"`
}

type AvailableTimeResponse struct {
	ID         uint      `json:"id"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	IsSelected bool      `json:"isSelected"`
} 