package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Permission struct {
	database.Model
	Type        enum.PermissionType `gorm:"not null;uniqueIndex"`
	Description string              `gorm:"type:text"`
	Category    enum.PermissionCategory
	Roles       []Role `gorm:"many2many:role_permissions;"`
}
