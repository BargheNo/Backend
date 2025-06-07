package service

import (
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
)

type NewsService interface {
	GetAllNewsStatuses() []newsdto.NewsStatusesResponse
	GetNews(request newsdto.GetNewsRequest) (newsdto.NewsResponse, error)
	GetNewsList(request newsdto.GetNewsListRequest) ([]newsdto.NewsResponse, error)
	CreateNews(request newsdto.CreateNewsRequest) (newsdto.NewsResponse, error)
	EditNews(request newsdto.EditNewsRequest) error
	UpdateNewsStatus(request newsdto.EditNewsStatusRequest) error
	DeleteNewsStatus(request newsdto.DeleteNewsRequest) error
	AddNewsMedia(request newsdto.AddNewsMediaRequest) (uint, error)
	DeleteNewsMedia(request newsdto.AccessMediaRequest) error
	GetNewsMedia(request newsdto.AccessMediaRequest) (string, error)
}
