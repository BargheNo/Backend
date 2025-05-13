package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository interface {
	FindNewsByTittle(db database.Database, title string) (*entity.News, bool)
	CreateNews(db database.Database, news *entity.News) error
}
