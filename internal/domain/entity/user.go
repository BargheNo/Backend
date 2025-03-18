package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type User struct {
	database.Model
	FirstName            string
	LastName             string
	Phone                string
	PhoneVerified        bool
	Password             string
	Email                string
	EmailVerified        bool
	NationalCode         string
	Address              []Address `gorm:"polymorphic:Owner;"`
	InstallationRequests []InstallationRequest
	Roles                []Role	`gorm:"many2many:user_roles;"`
}
