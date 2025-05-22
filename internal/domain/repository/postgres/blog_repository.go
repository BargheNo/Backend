package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BlogRepository interface {
	CreatePost(db database.Database, post *entity.Post) error
	UpdatePost(db database.Database, post *entity.Post) error
	GetCorporationPosts(db database.Database, corporationID uint, opts ...QueryModifier) []entity.Post
	FindPostByID(db database.Database, postID uint) (*entity.Post, bool)
}
