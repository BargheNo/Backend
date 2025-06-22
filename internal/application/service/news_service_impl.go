package service

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	postgresImpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type NewsService struct {
	constants      *bootstrap.Constants
	userService    port.UserService
	s3Storage      s3.S3Storage
	newsRepository postgres.NewsRepository
	db             database.Database
}

func NewNewsService(
	constants *bootstrap.Constants,
	userService port.UserService,
	s3Storage s3.S3Storage,
	newsRepository postgres.NewsRepository,
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

func (newsService *NewsService) mapToFilterStatuses(enumStatus uint) []enum.NewsStatus {
	statuses := enum.GetAllNewsStatus()
	for _, status := range statuses {
		if uint(status) == enumStatus {
			if status == enum.NewsStatusAll {
				return statuses
			}
			return []enum.NewsStatus{status}
		}
	}
	return statuses
}

func (newsService *NewsService) mapToOperationalStatuses(enumStatus uint) enum.NewsStatus {
	allowedStatuses := []enum.NewsStatus{enum.NewsStatusActive, enum.NewsStatusDraft}
	for _, status := range allowedStatuses {
		if uint(status) == enumStatus {
			return status
		}
	}
	return enum.NewsStatusDraft
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

func (newsService *NewsService) getNewsByID(newsID uint) (*entity.News, error) {
	news, err := newsService.newsRepository.FindNewsByID(newsService.db, newsID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return nil, notFoundError
	}
	return news, nil
}

func (newsService *NewsService) getNewsMedia(mediaID, newsID uint) (*entity.Media, error) {
	media, err := newsService.newsRepository.FindNewsMediaByID(newsService.db, mediaID, newsID, "news")
	if err != nil {
		return nil, err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		return nil, notFoundError
	}
	return media, nil
}

func (newsService *NewsService) GetAdminNews(newsID uint) (newsdto.AdminNewsResponse, error) {
	news, err := newsService.getNewsByID(newsID)
	if err != nil {
		return newsdto.AdminNewsResponse{}, err
	}

	coverImage := ""
	if news.CoverImage != "" {
		coverImage, err = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, news.CoverImage, 8*time.Hour)
		if err != nil {
			return newsdto.AdminNewsResponse{}, err
		}
	}

	author, err := newsService.userService.GetUserCredential(news.AuthorID)
	if err != nil {
		return newsdto.AdminNewsResponse{}, err
	}

	return newsdto.AdminNewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Description: news.Description,
		Status:      news.Status.String(),
		CoverImage:  coverImage,
		Author:      author,
	}, nil
}

func (newsService *NewsService) GetPublicNews(newsID uint) (newsdto.PublicNewsResponse, error) {
	news, err := newsService.getNewsByID(newsID)
	if err != nil {
		return newsdto.PublicNewsResponse{}, err
	}

	if news.Status != enum.NewsStatusActive {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.News}
		return newsdto.PublicNewsResponse{}, notFoundError
	}

	coverImage := ""
	if news.CoverImage != "" {
		coverImage, err = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, news.CoverImage, 8*time.Hour)
		if err != nil {
			return newsdto.PublicNewsResponse{}, err
		}
	}

	return newsdto.PublicNewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Description: news.Description,
		CoverImage:  coverImage,
	}, nil
}

func (newsService *NewsService) GetAdminNewsList(request newsdto.GetAdminNewsListRequest) ([]newsdto.AdminNewsResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	allowedStatuses := newsService.mapToFilterStatuses(request.Status)
	news, err := newsService.newsRepository.FindNewsByStatus(newsService.db, allowedStatuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	newsResponse := make([]newsdto.AdminNewsResponse, len(news))

	for i, eachNews := range news {
		coverImage := ""
		if eachNews.CoverImage != "" {
			coverImage, err = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, eachNews.CoverImage, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}

		author, err := newsService.userService.GetUserCredential(eachNews.AuthorID)
		if err != nil {
			return nil, err
		}

		newsResponse[i] = newsdto.AdminNewsResponse{
			ID:          eachNews.ID,
			Title:       eachNews.Title,
			Content:     eachNews.Content,
			Description: eachNews.Description,
			Status:      eachNews.Status.String(),
			CoverImage:  coverImage,
			Author:      author,
		}
	}
	return newsResponse, nil
}

func (newsService *NewsService) GetPublicNewsList(request newsdto.GetPublicNewsListRequest) ([]newsdto.PublicNewsResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	allowedStatuses := []enum.NewsStatus{enum.NewsStatusActive}
	news, err := newsService.newsRepository.FindNewsByStatus(newsService.db, allowedStatuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	newsResponse := make([]newsdto.PublicNewsResponse, len(news))

	for i, eachNews := range news {
		coverImage := ""
		if eachNews.CoverImage != "" {
			coverImage, err = newsService.s3Storage.GetPresignedURL(enum.NewsMedia, eachNews.CoverImage, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}

		newsResponse[i] = newsdto.PublicNewsResponse{
			ID:          eachNews.ID,
			Title:       eachNews.Title,
			Content:     eachNews.Content,
			Description: eachNews.Description,
			CoverImage:  coverImage,
		}
	}
	return newsResponse, nil
}

func (newsService *NewsService) checkDuplicateNews(title string) error {
	news, err := newsService.newsRepository.FindNewsByTittle(newsService.db, title)
	if err != nil {
		return err
	}
	if news != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(newsService.constants.Field.Name, newsService.constants.Tag.AlreadyExist)
		return conflictErrors
	}
	return nil
}

func (newsService *NewsService) CreateNews(request newsdto.CreateNewsRequest) (uint, error) {
	if err := newsService.userService.IsUserActive(request.AuthorID); err != nil {
		return 0, nil
	}

	if err := newsService.checkDuplicateNews(request.Title); err != nil {
		return 0, err
	}

	news := &entity.News{
		Title:       request.Title,
		Content:     request.Content,
		Description: request.Description,
		AuthorID:    request.AuthorID,
		Status:      request.Status,
	}
	err := newsService.db.WithTransaction(func(tx database.Database) error {
		if err := newsService.newsRepository.CreateNews(tx, news); err != nil {
			return err
		}

		if request.CoverImage != nil {
			news.CoverImage = newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
			if err := newsService.s3Storage.UploadObject(enum.NewsMedia, news.CoverImage, request.CoverImage); err != nil {
				return err
			}

			if err := newsService.newsRepository.UpdateNews(tx, news); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return news.ID, nil
}

func (newsService *NewsService) checkStatusConflict(newStatus, oldStatus enum.NewsStatus) error {
	var conflictErrors exception.ConflictErrors
	if newStatus == enum.NewsStatusActive && oldStatus == enum.NewsStatusActive {
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyActive)
		return conflictErrors
	}
	if newStatus == enum.NewsStatusDraft && oldStatus == enum.NewsStatusDraft {
		conflictErrors.Add(newsService.constants.Field.News, newsService.constants.Tag.AlreadyDraft)
		return conflictErrors
	}
	return nil
}

func (newsService *NewsService) EditNews(request newsdto.EditNewsRequest) error {
	if err := newsService.userService.IsUserActive(request.AuthorID); err != nil {
		return err
	}

	news, err := newsService.getNewsByID(request.NewsID)
	if err != nil {
		return err
	}

	if request.Title != nil {
		if err := newsService.checkDuplicateNews(*request.Title); err != nil {
			return err
		}
		news.Title = *request.Title
	}

	if request.Content != nil {
		news.Content = *request.Content
	}

	if request.Description != nil {
		news.Description = *request.Description
	}

	newStatus := newsService.mapToOperationalStatuses(request.Status)
	if err := newsService.checkStatusConflict(newStatus, news.Status); err != nil {
		return err
	}
	news.Status = newStatus

	if request.CoverImage != nil {
		news.CoverImage = newsService.constants.S3BucketPath.GetNewsCoverImagePath(news.ID, request.CoverImage.Filename)
		if err := newsService.s3Storage.UploadObject(enum.NewsMedia, news.CoverImage, request.CoverImage); err != nil {
			return err
		}
	}

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		return err
	}
	return nil
}

func (newsService *NewsService) UpdateNewsStatus(request newsdto.EditNewsStatusRequest) error {
	if err := newsService.userService.IsUserActive(request.AuthorID); err != nil {
		return err
	}

	news, err := newsService.getNewsByID(request.NewsID)
	if err != nil {
		return err
	}

	if err := newsService.checkStatusConflict(enum.NewsStatus(request.Status), news.Status); err != nil {
		return err
	}
	news.Status = enum.NewsStatus(request.Status)

	if err := newsService.newsRepository.UpdateNews(newsService.db, news); err != nil {
		return err
	}
	return nil
}

func (newsService *NewsService) DeleteNewsStatus(request newsdto.DeleteNewsRequest) error {
	err := newsService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}

	for _, newsID := range request.NewsIDs {
		news, err := newsService.newsRepository.FindNewsByID(newsService.db, newsID)
		if err != nil {
			return err
		}
		if news == nil {
			continue
		}

		if err := newsService.newsRepository.DeleteNews(newsService.db, newsID); err != nil {
			return err
		}
	}
	return nil
}

func (newsService *NewsService) AddNewsMedia(request newsdto.AddNewsMediaRequest) (uint, error) {
	if err := newsService.userService.IsUserActive(request.AuthorID); err != nil {
		return 0, err
	}

	if _, err := newsService.getNewsByID(request.NewsID); err != nil {
		return 0, err
	}

	mediaPath := newsService.constants.S3BucketPath.GetNewsMediaPath(request.NewsID, request.Media.Filename)
	if err := newsService.s3Storage.UploadObject(enum.NewsMedia, mediaPath, request.Media); err != nil {
		return 0, err
	}

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.NewsID,
		OwnerType: "news",
	}
	if err := newsService.newsRepository.CreateMedia(newsService.db, media); err != nil {
		return 0, err
	}
	return media.ID, nil
}

func (newsService *NewsService) DeleteNewsMedia(request newsdto.AccessMediaRequest) error {
	if err := newsService.userService.IsUserActive(request.AuthorID); err != nil {
		return err
	}

	if _, err := newsService.getNewsByID(request.NewsID); err != nil {
		return err
	}

	media, err := newsService.getNewsMedia(request.MediaID, request.NewsID)
	if err != nil {
		return err
	}

	err = newsService.db.WithTransaction(func(tx database.Database) error {
		if err := newsService.newsRepository.DeleteMedia(tx, request.MediaID); err != nil {
			return err
		}
		if err := newsService.s3Storage.DeleteObject(enum.NewsMedia, media.Path); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (newsService *NewsService) GetNewsMedia(request newsdto.AccessMediaRequest) (string, error) {
	news, err := newsService.getNewsByID(request.NewsID)
	if err != nil {
		return "", err
	}

	if request.UserType == enum.UserTypeGuest && news.Status == enum.NewsStatusDraft {
		notFoundError := exception.NotFoundError{Item: newsService.constants.Field.Media}
		return "", notFoundError
	}

	media, err := newsService.getNewsMedia(request.MediaID, request.NewsID)
	if err != nil {
		return "", err
	}

	presignedURL, err := newsService.s3Storage.GetPresignedURL(enum.NewsMedia, media.Path, 8*time.Hour)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}
