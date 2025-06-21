package service

import (
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
)

type NewsService interface {
	GetAllNewsStatuses() []newsdto.NewsStatusesResponse
	GetAdminNews(newsID uint) (newsdto.AdminNewsResponse, error)
	GetPublicNews(newsID uint) (newsdto.PublicNewsResponse, error)
	GetAdminNewsList(request newsdto.GetAdminNewsListRequest) ([]newsdto.AdminNewsResponse, error)
	GetPublicNewsList(request newsdto.GetPublicNewsListRequest) ([]newsdto.PublicNewsResponse, error)
	CreateNews(request newsdto.CreateNewsRequest) (uint, error)
	EditNews(request newsdto.EditNewsRequest) error
	UpdateNewsStatus(request newsdto.EditNewsStatusRequest) error
	DeleteNewsStatus(request newsdto.DeleteNewsRequest) error
	AddNewsMedia(request newsdto.AddNewsMediaRequest) (uint, error)
	DeleteNewsMedia(request newsdto.AccessMediaRequest) error
	GetNewsMedia(request newsdto.AccessMediaRequest) (string, error)
}
