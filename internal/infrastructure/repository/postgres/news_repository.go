package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository struct {
}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{}
}

func (repo *UserRepository) CreateNews(db database.Database, news *entity.News) error {
	return db.GetDB().Create(&news).Error
}
