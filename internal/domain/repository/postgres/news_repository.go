package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsRepository interface {
	CreateNews(db database.Database, news *entity.News) error
}
