package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Permission struct {
	database.Model
	Name        string
	Description string
	Roles       []Role `gorm:"many2many:role_permissions;"`
}
