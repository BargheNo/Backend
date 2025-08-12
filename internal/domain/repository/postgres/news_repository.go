package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository interface {
	FindNewsByID(db database.Database, newsID uint) (*entity.News, error)
	FindNewsByTittle(db database.Database, title string) (*entity.News, error)
	FindNewsByStatus(db database.Database, statuses []enum.NewsStatus, options *QueryOptions) ([]*entity.News, error)
	CountNewsByStatus(db database.Database, statuses []enum.NewsStatus) (int64, error)
	FindNewsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.News, error)
	CountNewsByQuery(db database.Database, query string) (int64, error)
	UpdateNews(db database.Database, news *entity.News) error
	CreateNews(db database.Database, news *entity.News) error
	DeleteNews(db database.Database, newsID uint) error
	FindNewsMediaByID(db database.Database, mediaID, newsID uint, ownerType string) (*entity.Media, error)
	CreateMedia(db database.Database, media *entity.Media) error
	DeleteMedia(db database.Database, mediaID uint) error
	FindNewsLikeByUser(db database.Database, userID, newsID uint) (*entity.Like, error)
	CreateLike(db database.Database, userID uint, newsID uint) error
	DeleteLike(db database.Database, like *entity.Like) error
}
