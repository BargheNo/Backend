package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type User struct {
	database.Model
	FirstName     string          `gorm:"type:varchar(50);index:idx_user_name"`
	LastName      string          `gorm:"type:varchar(50);index:idx_user_name"`
	Phone         string          `gorm:"type:varchar(20);uniqueIndex"`
	PhoneVerified bool            `gorm:"default:false"`
	Password      string          `gorm:"type:varchar(255);not null"`
	Email         string          `gorm:"type:varchar(100);Index"`
	EmailVerified bool            `gorm:"default:false"`
	NationalCode  string          `gorm:"type:varchar(20);Index"`
	Status        enum.UserStatus `gorm:"index"`
	Addresses     []Address       `gorm:"polymorphic:Owner;polymorphicValue:users"`
	Roles         []Role          `gorm:"many2many:user_roles;"`
}
