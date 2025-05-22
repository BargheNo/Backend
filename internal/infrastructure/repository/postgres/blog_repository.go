package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
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

func (repo *BlogRepository) GetCorporationPosts(db database.Database, corporationID uint) ([]entity.Post, error) {
	var posts []entity.Post
	if err := db.GetDB().Where("corporation_id = ?", corporationID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
