package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Role struct {
	database.Model
	Name        string       `gorm:"type:varchar(50);uniqueIndex"`
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
