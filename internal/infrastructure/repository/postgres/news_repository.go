package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (repo *NewsRepository) FindNewsByStatus(db database.Database, statues []uint, opts ...repository.QueryModifier) []*entity.News {
	var news []*entity.News
	query := db.GetDB().Where("status IN ?", statues)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&news)
	if result.Error != nil {
		panic(result.Error)
	}
	return news
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

func (repo *NewsRepository) GetMediaByID(db database.Database, mediaID uint) (*entity.Media, bool) {
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

func (repo *NewsRepository) AddMedia(db database.Database, media *entity.Media) error {
	return db.GetDB().Create(&media).Error
}

func (repo *NewsRepository) DeleteMedia(db database.Database, mediaID uint) error {
	return db.GetDB().Delete(&entity.Media{}, mediaID).Error
}
