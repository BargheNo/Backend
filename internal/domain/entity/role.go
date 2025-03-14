package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Role struct {
	database.Model
	Name        string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
