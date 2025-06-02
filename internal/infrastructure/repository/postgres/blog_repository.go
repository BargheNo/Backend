package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BlogRepository struct{}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{}
}

func (repo *BlogRepository) CreatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Create(post).Error
}

func (repo *BlogRepository) UpdatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Save(post).Error
}

func (repo *BlogRepository) GetCorporationPostsByStatus(db database.Database, corporationID uint, statuses []uint, opts ...repository.QueryModifier) []entity.Post {
	var posts []entity.Post
	query := db.GetDB().Where("corporation_id = ?", corporationID).Where("status IN (?)", statuses)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&posts)
	if result.Error != nil {
		panic(result.Error)
	}
	return posts
}

func (repo *BlogRepository) FindPostByID(db database.Database, postID uint) (*entity.Post, bool) {
	var post entity.Post
	result := db.GetDB().Where("id = ?", postID).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &post, true
}

func (repo *BlogRepository) GetMediaByID(db database.Database, mediaID uint) (*entity.Media, bool) {
	var media entity.Media
	result := db.GetDB().First(&media, mediaID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &media, true
}

func (repo *BlogRepository) AddMedia(db database.Database, media *entity.Media) error {
	return db.GetDB().Create(&media).Error
}

func (repo *BlogRepository) DeleteMedia(db database.Database, mediaID uint) error {
	return db.GetDB().Delete(&entity.Media{}, mediaID).Error
}

func (repo *BlogRepository) DeletePost(db database.Database, postID uint) error {
	return db.GetDB().Delete(&entity.Post{}, postID).Error
}

func (repo *BlogRepository) CreateLike(db database.Database, like *entity.Like) error {
	return db.GetDB().Create(&like).Error
}

func (repo *BlogRepository) FindLikeByUserAndOwner(db database.Database, userID, ownerID uint, ownerType string) (*entity.Like, bool) {
	var like entity.Like
	result := db.GetDB().Where("user_id = ? AND owner_id = ? AND owner_type = ?", userID, ownerID, ownerType).First(&like)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &like, true
}

func (repo *BlogRepository) DeleteLike(db database.Database, likeID uint) error {
	return db.GetDB().Delete(&entity.Like{}, likeID).Error
}

func (repo *BlogRepository) GetPostsByStatus(db database.Database, statuses []uint, opts ...repository.QueryModifier) []entity.Post {
	var posts []entity.Post
	query := db.GetDB().Where("status IN (?)", statuses)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&posts)
	if result.Error != nil {
		panic(result.Error)
	}
	return posts
}

func (repo *BlogRepository) GetLikeCountByOwner(db database.Database, ownerID uint, ownerType string) uint {
	var count int64
	db.GetDB().Model(&entity.Like{}).Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Count(&count)
	return uint(count)
}
