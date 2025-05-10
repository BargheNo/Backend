package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NewsService struct {
	constants      *bootstrap.Constants
	userService    service.UserService
	newsRepository repository.NewsRepository
	db             database.Database
}

func NewNewsService(
	constants *bootstrap.Constants,
	userService service.UserService,
	newsRepository repository.NewsRepository,
	db database.Database,
) *NewsService {
	return &NewsService{
		constants:      constants,
		userService:    userService,
		newsRepository: newsRepository,
		db:             db,
	}
}

func (newsService *NewsService) CreateNews(request newsdto.CreateNewsRequest) {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}
	news := &entity.News{
		Title:    request.Title,
		Content:  request.Content,
		AuthorID: request.AuthorID,
		Status:   request.Status,
	}
	if err := newsService.newsRepository.CreateNews(newsService.db, news); err != nil {
		panic(err)
	}
}
