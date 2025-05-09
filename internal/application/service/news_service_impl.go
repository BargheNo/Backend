package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsService struct {
	constants      *bootstrap.Constants
	newsRepository repository.NewsRepository
	db             database.Database
}

func NewNewsService(
	constants *bootstrap.Constants,
	newsRepository repository.NewsRepository,
	db database.Database,
) *NewsService {
	return &NewsService{
		constants:      constants,
		newsRepository: newsRepository,
		db:             db,
	}
}
