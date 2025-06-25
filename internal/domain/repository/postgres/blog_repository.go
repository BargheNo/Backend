package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BlogRepository interface {
	CreateLike(db database.Database, like *entity.Like) error
	CreateMedia(db database.Database, media *entity.Media) error
	CreatePost(db database.Database, post *entity.Post) error
	DeleteLike(db database.Database, likeID uint) error
	DeleteMedia(db database.Database, mediaID uint) error
	DeletePost(db database.Database, postID uint) error
	FindCorporationPost(db database.Database, postID uint, corporationID uint) (*entity.Post, error)
	FindCorporationPostByTitle(db database.Database, corporationID uint, title string) (*entity.Post, error)
	FindCorporationPostsByStatus(db database.Database, corporationID uint, statuses []enum.PostStatus, opts ...QueryModifier) ([]entity.Post, error)
	FindLikeByUserAndOwner(db database.Database, userID uint, ownerID uint, ownerType string) (*entity.Like, error)
	FindLikeCountByOwner(db database.Database, ownerID uint, ownerType string) (uint, error)
	FindPostByID(db database.Database, postID uint) (*entity.Post, error)
	FindPostsByStatus(db database.Database, statuses []enum.PostStatus, opts ...QueryModifier) ([]entity.Post, error)
	FindPostMediaByID(db database.Database, mediaID, postID uint, ownerType string) (*entity.Media, error)
	UpdatePost(db database.Database, post *entity.Post) error
}
