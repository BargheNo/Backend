package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository interface {
	FindNewsByID(db database.Database, newsID uint) (*entity.News, bool)
	FindNewsByTittle(db database.Database, title string) (*entity.News, bool)
	FindNewsByStatus(db database.Database, statues []uint, opts ...QueryModifier) []*entity.News
	UpdateNews(db database.Database, news *entity.News) error
	CreateNews(db database.Database, news *entity.News) error
	DeleteNews(db database.Database, newsID uint) error
	GetMediaByID(db database.Database, mediaID uint) (*entity.Media, bool)
	AddMedia(db database.Database, media *entity.Media) error
	DeleteMedia(db database.Database, mediaID uint) error
}
