package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type News struct {
	database.Model
	Title    string `json:"title"`
	Content  string `json:"content_html"`
	AuthorID uint   `gorm:"not null;index"`
	Author   User   `gorm:"foreignKey:AuthorID"`
	Status   enum.NewsStatus
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

// type Image struct {
// 	ID       string
// DocID	 uint
// Doc 		 Document
// 	S3Key    string
// 	Caption  string
// 	Position int
// }
