package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (repo *NewsRepository) FindNewsByStatus(db database.Database, statuses []enum.NewsStatus, options *postgres.QueryOptions) ([]*entity.News, error) {
	var news []*entity.News
	query := db.GetDB().Where("status IN ?", statuses)
	query = applyQueryOptions(query, options)
	result := query.Find(&news)
	if result.Error != nil {
		return nil, result.Error
	}
	return news, nil
}

func (repo *NewsRepository) CountNewsByStatus(db database.Database, statuses []enum.NewsStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.News{}).
		Where("status IN ?", statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *NewsRepository) FindNewsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.News, error) {
	var news []*entity.News
	result := db.GetDB().
		Joins("LEFT JOIN users AS authors ON news.author_id = authors.id").
		Where("title ILIKE ? OR description ILIKE ? OR authors.first_name ILIKE ? OR authors.last_name ILIKE ? OR authors.email ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	result = applyQueryOptions(result, options)
	result = result.Find(&news)
	if result.Error != nil {
		return nil, result.Error
	}
	return news, nil
}

func (repo *NewsRepository) CountNewsByQuery(db database.Database, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.News{}).
		Joins("LEFT JOIN users AS authors ON news.author_id = authors.id").
		Where("title ILIKE ? OR description ILIKE ? OR authors.first_name ILIKE ? OR authors.last_name ILIKE ? OR authors.email ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
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

func (repo *NewsRepository) FindNewsLikeByUser(db database.Database, userID, newsID uint) (*entity.Like, error) {
	var like entity.Like
	result := db.GetDB().Where("user_id = ? AND owner_id = ? AND owner_type = ?", userID, newsID, "news").First(&like)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &like, nil
}

func (repo *NewsRepository) CreateLike(db database.Database, userID uint, newsID uint) error {
	like := &entity.Like{
		UserID:    userID,
		OwnerID:   newsID,
		OwnerType: "news",
	}
	return db.GetDB().Create(&like).Error
}

func (repo *NewsRepository) DeleteLike(db database.Database, like *entity.Like) error {
	return db.GetDB().Unscoped().Delete(like).Error
}
