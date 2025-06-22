package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type NewsRepository struct {
}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{}
}

func (repo *NewsRepository) FindNewsByID(db database.Database, newsID uint) (*entity.News, error) {
	var news entity.News
	result := db.GetDB().First(&news, newsID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &news, nil
}

func (repo *NewsRepository) FindNewsByTittle(db database.Database, title string) (*entity.News, error) {
	var news entity.News
	result := db.GetDB().Where("title = ?", title).First(&news)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &news, nil
}

func (repo *NewsRepository) FindNewsByStatus(db database.Database, statuses []enum.NewsStatus, opts ...repository.QueryModifier) ([]*entity.News, error) {
	var news []*entity.News
	query := db.GetDB().Where("status IN ?", statuses)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&news)
	if result.Error != nil {
		return nil, result.Error
	}
	return news, nil
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

func (repo *NewsRepository) FindNewsMediaByID(db database.Database, mediaID, newsID uint, ownerType string) (*entity.Media, error) {
	var media entity.Media
	result := db.GetDB().Where("id = ? AND owner_id = ? AND owner_type = ?", mediaID, newsID, ownerType).First(&media)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &media, nil
}

func (repo *NewsRepository) CreateMedia(db database.Database, media *entity.Media) error {
	return db.GetDB().Create(&media).Error
}

func (repo *NewsRepository) DeleteMedia(db database.Database, mediaID uint) error {
	return db.GetDB().Delete(&entity.Media{}, mediaID).Error
}
