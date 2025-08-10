package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BlogRepository struct{}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{}
}

func (repo *BlogRepository) FindPostByID(db database.Database, postID uint) (*entity.Post, error) {
	var post entity.Post
	result := db.GetDB().First(&post, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &post, nil
}

func (repo *BlogRepository) FindCorporationPost(db database.Database, postID, corporationID uint) (*entity.Post, error) {
	var post entity.Post
	result := db.GetDB().Where("id = ? AND corporation_id = ?", postID, corporationID).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &post, nil
}

func (repo *BlogRepository) FindCorporationPostByTitle(db database.Database, corporationID uint, title string) (*entity.Post, error) {
	var post entity.Post
	result := db.GetDB().Where("corporation_id = ? AND title = ?", corporationID, title).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &post, nil
}

func (repo *BlogRepository) FindCorporationPostsByStatus(db database.Database, corporationID uint, statuses []enum.PostStatus, options *postgres.QueryOptions) ([]entity.Post, error) {
	var posts []entity.Post
	query := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, statuses)
	query = applyQueryOptions(query, options)
	result := query.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (repo *BlogRepository) CountCorporationPostsByStatus(db database.Database, corporationID uint, statuses []enum.PostStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Post{}).
		Where("corporation_id = ? AND status IN ?", corporationID, statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *BlogRepository) CreatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Create(post).Error
}

func (repo *BlogRepository) UpdatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Save(post).Error
}

func (repo *BlogRepository) DeletePost(db database.Database, postID uint) error {
	return db.GetDB().Delete(&entity.Post{}, postID).Error
}

func (repo *BlogRepository) FindPostMediaByID(db database.Database, mediaID, postID uint, ownerType string) (*entity.Media, error) {
	var media entity.Media
	result := db.GetDB().Where("id = ? AND owner_id = ? AND owner_type = ?", mediaID, postID, ownerType).First(&media)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &media, nil
}

func (repo *BlogRepository) CreateMedia(db database.Database, media *entity.Media) error {
	return db.GetDB().Create(&media).Error
}

func (repo *BlogRepository) DeleteMedia(db database.Database, mediaID uint) error {
	return db.GetDB().Delete(&entity.Media{}, mediaID).Error
}

func (repo *BlogRepository) FindLikeByUserAndBlogID(db database.Database, userID, ownerID uint) (*entity.Like, error) {
	var like entity.Like
	result := db.GetDB().Where("user_id = ? AND owner_id = ? AND owner_type = ?", userID, ownerID, "blog").First(&like)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &like, nil
}

func (repo *BlogRepository) FindPostsByStatus(db database.Database, statuses []enum.PostStatus, options *postgres.QueryOptions) ([]entity.Post, error) {
	var posts []entity.Post
	query := db.GetDB().Where("status IN ?", statuses)
	query = applyQueryOptions(query, options)
	result := query.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (repo *BlogRepository) CountPostsByStatus(db database.Database, statuses []enum.PostStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Post{}).
		Where("status IN ?", statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *BlogRepository) CreateLike(db database.Database, userID, blogID uint) error {
	like := &entity.Like{
		UserID:    userID,
		OwnerID:   blogID,
		OwnerType: "blog",
	}
	return db.GetDB().Create(&like).Error
}

func (repo *BlogRepository) DeleteLike(db database.Database, likeID uint) error {
	return db.GetDB().Delete(&entity.Like{}, likeID).Error
}
