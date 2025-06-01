package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Post struct {
	database.Model
	Title         string      `json:"title"`
	CoverImage    string      `gorm:"type:varchar(255);default:null"`
	Content       string      `json:"content_html"`
	AuthorID      uint        `gorm:"not null;index"`
	Author        User        `gorm:"foreignKey:AuthorID"`
	CorporationID uint        `gorm:"not null;index"`
	Corporation   Corporation `gorm:"foreignKey:CorporationID"`
	Media         []Media     `gorm:"polymorphic:Owner;polymorphicValue:posts"`
	Status        enum.PostStatus
	Likes         []Like `gorm:"polymorphic:Owner;polymorphicValue:posts"`
}
