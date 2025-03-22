package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type ContactInformation struct {
	database.Model
	Phone         string
	Email         string
	Eitaa         string
	Bale          string
	Website       string
	WhatsApp      string
	Instagram     string
	Linkedin      string
	Telegram      string
	CorporationID uint
}
