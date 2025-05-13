package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
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

func (newsService *NewsService) CreateNews(request newsdto.CreateNewsRequest) newsdto.NewsResponse {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}
	_, exist := newsService.newsRepository.FindNewsByTittle(newsService.db, request.Title)
	if exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.Name, newsService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
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
	newsResponse := newsdto.NewsResponse{
		ID:      news.ID,
		Title:   news.Title,
		Content: news.Content,
		Status:  news.Status,
	}
	return newsResponse
}

func (newsService *NewsService) EditNews(request newsdto.EditNewsRequest) {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}
	news, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}

	if request.Title != nil {
		news.Title = *request.Title
	}
	if request.Content != nil {
		news.Content = *request.Content
	}
	news.Status = enum.NewsStatus(request.Status)

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		panic(err)
	}
}

func (newsService *NewsService) UpdateNewsStatus(request newsdto.EditNewsStatusRequest) {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}

	news, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}

	if enum.NewsStatus(request.Status) == enum.NewsStatusActive && news.Status == enum.NewsStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyActive)
		panic(conflictErrors)
	}
	if enum.NewsStatus(request.Status) == enum.NewsStatusDraft && news.Status == enum.NewsStatusDraft {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyDraft)
		panic(conflictErrors)
	}

	news.Status = enum.NewsStatus(request.Status)
	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		panic(err)
	}
}
