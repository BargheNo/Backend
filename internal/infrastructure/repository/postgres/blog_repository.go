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
