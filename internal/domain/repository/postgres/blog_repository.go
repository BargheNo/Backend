package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BlogRepository interface {
	CreatePost(db database.Database, post *entity.Post) error
	UpdatePost(db database.Database, post *entity.Post) error
	GetCorporationPostsByStatus(db database.Database, corporationID uint, statuses []uint, opts ...QueryModifier) ([]entity.Post, error)
	GetPostsByStatus(db database.Database, statuses []uint, opts ...QueryModifier) ([]entity.Post, error)
	FindPostByID(db database.Database, postID uint) (*entity.Post, error)
	GetMediaByID(db database.Database, mediaID uint) (*entity.Media, error)
	AddMedia(db database.Database, media *entity.Media) error
	DeleteMedia(db database.Database, mediaID uint) error
	DeletePost(db database.Database, postID uint) error
	CreateLike(db database.Database, like *entity.Like) error
	FindLikeByUserAndOwner(db database.Database, userID, ownerID uint, ownerType string) (*entity.Like, error)
	DeleteLike(db database.Database, likeID uint) error
	GetLikeCountByOwner(db database.Database, ownerID uint, ownerType string) uint
}
