package serviceimpl

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type NewsService struct {
	constants      *bootstrap.Constants
	userService    service.UserService
	s3Storage      s3.S3Storage
	newsRepository repository.NewsRepository
	db             database.Database
}

func NewNewsService(
	constants *bootstrap.Constants,
	userService service.UserService,
	s3Storage s3.S3Storage,
	newsRepository repository.NewsRepository,
	db database.Database,
) *NewsService {
	return &NewsService{
		constants:      constants,
		userService:    userService,
		s3Storage:      s3Storage,
		newsRepository: newsRepository,
		db:             db,
	}
}

func (newsService *NewsService) GetAllNewsStatuses() []newsdto.NewsStatusesResponse {
	allowedStatuses := []enum.NewsStatus{
		enum.NewsStatusActive,
		enum.NewsStatusDraft,
	}

	statuses := make([]newsdto.NewsStatusesResponse, len(allowedStatuses))
	for i, status := range allowedStatuses {
		statuses[i] = newsdto.NewsStatusesResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return statuses
}

func (newsService *NewsService) GetNews(request newsdto.GetNewsRequest) newsdto.NewsResponse {
	news, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}
	if request.UserType == enum.UserTypeGuest && news.Status == enum.NewsStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}
	coverImage := ""
	if news.CoverImage != "" {
		coverImage = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, news.CoverImage, 8*time.Hour)
	}
	return newsdto.NewsResponse{
		ID:         news.ID,
		Title:      news.Title,
		Content:    news.Content,
		Status:     news.Status,
		CoverImage: coverImage,
	}
}

func (newsService *NewsService) GetNewsList(request newsdto.GetNewsListRequest) []newsdto.NewsResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	news := newsService.newsRepository.FindNewsByStatus(newsService.db, request.Statuses, paginationModifier, sortingModifier)
	newsResponse := make([]newsdto.NewsResponse, len(news))
	for i, eachNews := range news {
		coverImage := ""
		if eachNews.CoverImage != "" {
			coverImage = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, eachNews.CoverImage, 8*time.Hour)
		}
		newsResponse[i] = newsdto.NewsResponse{
			ID:         eachNews.ID,
			Title:      eachNews.Title,
			Content:    eachNews.Content,
			Status:     eachNews.Status,
			CoverImage: coverImage,
		}
	}
	return newsResponse
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

	if request.CoverImage != nil {
		mediaPath := newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
		newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.CoverImage)
		news.CoverImage = mediaPath
	}

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
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

	if request.CoverImage != nil {
		mediaPath := newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
		newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.CoverImage)
		news.CoverImage = mediaPath
	}

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

func (newsService *NewsService) DeleteNewsStatus(request newsdto.DeleteNewsRequest) {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}

	for _, newsID := range request.NewsIDs {
		_, exist := newsService.newsRepository.FindNewsByID(newsService.db, newsID)
		if !exist {
			continue
		}
		newsService.newsRepository.DeleteNews(newsService.db, newsID)
	}
}

func (newsService *NewsService) AddNewsMedia(request newsdto.AddNewsMediaRequest) uint {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}
	_, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}
	mediaPath := newsService.constants.S3BucketPath.GetNewsMediaPath(request.NewsID, request.Media.Filename)
	newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.Media)

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.NewsID,
		OwnerType: "news",
	}
	if err := newsService.newsRepository.AddMedia(newsService.db, media); err != nil {
		panic(err)
	}
	return media.ID
}

func (newsService *NewsService) DeleteNewsMedia(request newsdto.AccessMediaRequest) {
	ok := newsService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}

	_, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}

	media, exist := newsService.newsRepository.GetMediaByID(newsService.db, request.MediaID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		panic(notFoundError)
	}

	if media.OwnerID != request.NewsID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.Media,
		}
		panic(forbiddenError)
	}

	if err := newsService.s3Storage.DeleteObject(enum.NewsMedia, media.Path); err != nil {
		panic(err)
	}
	if err := newsService.newsRepository.DeleteMedia(newsService.db, request.MediaID); err != nil {
		panic(err)
	}
}

func (newsService *NewsService) GetNewsMedia(request newsdto.AccessMediaRequest) string {
	news, exist := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		panic(notFoundError)
	}
	if request.UserType == enum.UserTypeGuest && news.Status == enum.NewsStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		panic(forbiddenError)
	}

	media, exist := newsService.newsRepository.GetMediaByID(newsService.db, request.MediaID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		panic(notFoundError)
	}

	if media.OwnerID != request.NewsID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.Media,
		}
		panic(forbiddenError)
	}

	return newsService.s3Storage.GetPresignedURL(enum.NewsMedia, media.Path, 8*time.Hour)
}
