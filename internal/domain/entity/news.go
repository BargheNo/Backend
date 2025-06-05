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
	CoverImage  string  `gorm:"type:varchar(255);default:null"`
	Media       []Media `gorm:"polymorphic:Owner;polymorphicValue:news"`
	Status      enum.NewsStatus
	Likes       []Like `gorm:"polymorphic:Owner;polymorphicValue:news"`
}

// TODO: Better to update to this version
// type Document struct {
// 	ID       string
// 	Title    string
// 	Content  string
// 	Images   []Image
//  WriterID uint
//  Writer 	 User
//  Categories  []Category  `gorm:"many2many:news_categories;constraint:OnDelete:CASCADE"`
//  like count
//  subscribers
// }
