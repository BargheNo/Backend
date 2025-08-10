package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type News struct {
	database.Model
	Title       string  `json:"title"`
	Content     string  `json:"content_html"`
	Description string  `json:"description"`
	AuthorID    uint    `gorm:"not null;index"`
	Author      User    `gorm:"foreignKey:AuthorID"`
	CoverImage  string  `gorm:"type:text;default:null"`
	Media       []Media `gorm:"polymorphic:Owner;polymorphicValue:news"`
	Likes       []Like  `gorm:"polymorphic:Owner;polymorphicValue:news"`
	LikeCount   int     `gorm:"not null;default:0"`
	Status      enum.NewsStatus
}
