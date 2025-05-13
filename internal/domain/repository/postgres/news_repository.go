package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository interface {
	FindNewsByID(db database.Database, newsID uint) (*entity.News, bool)
	FindNewsByTittle(db database.Database, title string) (*entity.News, bool)
	UpdateNews(db database.Database, news *entity.News) error
	CreateNews(db database.Database, news *entity.News) error
}
