package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type User struct {
	database.Model
	FirstName     string
	LastName      string
	Phone         string
	Password      string
	OTP           string
	ResendAttempt uint
	Email         string
	NationalCode  string
	Address       []Address `gorm:"polymorphic:Owner;"`
	Roles         []Role    `gorm:"many2many:user_roles;"`
}
