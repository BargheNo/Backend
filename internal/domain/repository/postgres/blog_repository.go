package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BlogRepository interface {
	CreatePost(db database.Database, post *entity.Post) error
	UpdatePost(db database.Database, post *entity.Post) error
	GetCorporationPostsByStatus(db database.Database, corporationID uint, statuses []uint, opts ...QueryModifier) []entity.Post
	FindPostByID(db database.Database, postID uint) (*entity.Post, bool)
	GetMediaByID(db database.Database, mediaID uint) (*entity.Media, bool)
	AddMedia(db database.Database, media *entity.Media) error
	DeleteMedia(db database.Database, mediaID uint) error
	DeletePost(db database.Database, postID uint) error
	CreateLike(db database.Database, like *entity.Like) error
	FindLikeByUserAndOwner(db database.Database, userID, ownerID uint, ownerType string) (*entity.Like, bool)
	DeleteLike(db database.Database, likeID uint) error
}
