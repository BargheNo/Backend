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

func (newsService *NewsService) GetNews(request newsdto.GetNewsRequest) (newsdto.NewsResponse, error) {
	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return newsdto.NewsResponse{}, err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return newsdto.NewsResponse{}, notFoundError
	}
	if request.UserType == enum.UserTypeGuest && news.Status == enum.NewsStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return newsdto.NewsResponse{}, forbiddenError
	}
	coverImage := ""
	if news.CoverImage != "" {
		coverImage = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, news.CoverImage, 8*time.Hour)
	}
	return newsdto.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Description: news.Description,
		Status:      news.Status,
		CoverImage:  coverImage,
	}, nil
}

func (newsService *NewsService) GetNewsList(request newsdto.GetNewsListRequest) ([]newsdto.NewsResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	news, err := newsService.newsRepository.FindNewsByStatus(newsService.db, request.Statuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	newsResponse := make([]newsdto.NewsResponse, len(news))
	for i, eachNews := range news {
		coverImage := ""
		if eachNews.CoverImage != "" {
			coverImage = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, eachNews.CoverImage, 8*time.Hour)
		}
		newsResponse[i] = newsdto.NewsResponse{
			ID:          eachNews.ID,
			Title:       eachNews.Title,
			Content:     eachNews.Content,
			Description: eachNews.Description,
			Status:      eachNews.Status,
			CoverImage:  coverImage,
		}
	}
	return newsResponse, nil
}

func (newsService *NewsService) CreateNews(request newsdto.CreateNewsRequest) (newsdto.NewsResponse, error) {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return newsdto.NewsResponse{}, nil
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return newsdto.NewsResponse{}, forbiddenError
	}
	news, err := newsService.newsRepository.FindNewsByTittle(newsService.db, request.Title)
	if err != nil {
		return newsdto.NewsResponse{}, err
	}
	if news != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.Name, newsService.constants.Tag.AlreadyExist)
		return newsdto.NewsResponse{}, conflictErrors
	}
	news = &entity.News{
		Title:       request.Title,
		Content:     request.Content,
		Description: request.Description,
		AuthorID:    request.AuthorID,
		Status:      request.Status,
	}
	if err := newsService.newsRepository.CreateNews(newsService.db, news); err != nil {
		return newsdto.NewsResponse{}, err
	}

	if request.CoverImage != nil {
		mediaPath := newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
		newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.CoverImage)
		news.CoverImage = mediaPath
	}

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		return newsdto.NewsResponse{}, err
	}

	newsResponse := newsdto.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Description: news.Description,
		Status:      news.Status,
	}
	return newsResponse, nil
}

func (newsService *NewsService) EditNews(request newsdto.EditNewsRequest) error {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return forbiddenError
	}
	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return notFoundError
	}

	if request.Title != nil {
		news.Title = *request.Title
	}
	if request.Content != nil {
		news.Content = *request.Content
	}
	if request.Description != nil {
		news.Description = *request.Description
	}
	news.Status = enum.NewsStatus(request.Status)

	if request.CoverImage != nil {
		mediaPath := newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
		newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.CoverImage)
		news.CoverImage = mediaPath
	}

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		return err
	}
	return nil
}

func (newsService *NewsService) UpdateNewsStatus(request newsdto.EditNewsStatusRequest) error {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return forbiddenError
	}

	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return notFoundError
	}

	if enum.NewsStatus(request.Status) == enum.NewsStatusActive && news.Status == enum.NewsStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyActive)
		return conflictErrors
	}
	if enum.NewsStatus(request.Status) == enum.NewsStatusDraft && news.Status == enum.NewsStatusDraft {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyDraft)
		return conflictErrors
	}

	news.Status = enum.NewsStatus(request.Status)
	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		return err
	}
	return nil
}

func (newsService *NewsService) DeleteNewsStatus(request newsdto.DeleteNewsRequest) error {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return forbiddenError
	}

	for _, newsID := range request.NewsIDs {
		news, err := newsService.newsRepository.FindNewsByID(newsService.db, newsID)
		if err != nil {
			return err
		}
		if news == nil {
			continue
		}
		newsService.newsRepository.DeleteNews(newsService.db, newsID)
	}
	return nil
}

func (newsService *NewsService) AddNewsMedia(request newsdto.AddNewsMediaRequest) (uint, error) {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return 0, err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return 0, forbiddenError
	}
	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return 0, err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return 0, notFoundError
	}
	mediaPath := newsService.constants.S3BucketPath.GetNewsMediaPath(request.NewsID, request.Media.Filename)
	newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.Media)

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.NewsID,
		OwnerType: "news",
	}
	if err := newsService.newsRepository.AddMedia(newsService.db, media); err != nil {
		return 0, err
	}
	return media.ID, nil
}

func (newsService *NewsService) DeleteNewsMedia(request newsdto.AccessMediaRequest) error {
	ok, err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return forbiddenError
	}

	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return notFoundError
	}

	media, err := newsService.newsRepository.GetMediaByID(newsService.db, request.MediaID)
	if err != nil {
		return err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		return notFoundError
	}

	if media.OwnerID != request.NewsID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.Media,
		}
		return forbiddenError
	}

	if err := newsService.s3Storage.DeleteObject(enum.NewsMedia, media.Path); err != nil {
		return err
	}
	if err := newsService.newsRepository.DeleteMedia(newsService.db, request.MediaID); err != nil {
		return err
	}
	return nil
}

func (newsService *NewsService) GetNewsMedia(request newsdto.AccessMediaRequest) (string, error) {
	news, err := newsService.newsRepository.FindNewsByID(newsService.db, request.NewsID)
	if err != nil {
		return "", err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return "", notFoundError
	}
	if request.UserType == enum.UserTypeGuest && news.Status == enum.NewsStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.News,
		}
		return "", forbiddenError
	}

	media, err := newsService.newsRepository.GetMediaByID(newsService.db, request.MediaID)
	if err != nil {
		return "", err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		return "", notFoundError
	}

	if media.OwnerID != request.NewsID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: newsService.constants.Field.Media,
		}
		return "", forbiddenError
	}

	return newsService.s3Storage.GetPresignedURL(enum.NewsMedia, media.Path, 8*time.Hour), nil
}
