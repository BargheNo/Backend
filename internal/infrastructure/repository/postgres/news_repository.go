package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type NewsRepository struct {
}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{}
}

func (repo *NewsRepository) FindNewsByID(db database.Database, newsID uint) (*entity.News, bool) {
	var news entity.News
	result := db.GetDB().First(&news, newsID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *NewsRepository) FindNewsByTittle(db database.Database, title string) (*entity.News, bool) {
	var news entity.News
	result := db.GetDB().Where("title = ?", title).First(&news)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &news, true
}

func (repo *NewsRepository) UpdateNews(db database.Database, news *entity.News) error {
	return db.GetDB().Save(&news).Error
}

func (repo *NewsRepository) CreateNews(db database.Database, news *entity.News) error {
	return db.GetDB().Create(&news).Error
}

func (repo *NewsRepository) DeleteNews(db database.Database, newsID uint) error {
	return db.GetDB().Delete(&entity.News{}, newsID).Error
}
