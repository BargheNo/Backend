package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type ContactInformation struct {
	database.Model
	Address   []Address `gorm:"polymorphic:Owner;"`
	WhatsApp  string
	Instagram string
	Telegram  string
	Phone     string
	Email     string
	Eitaa     string
	Bale      string
	Website   string
}
