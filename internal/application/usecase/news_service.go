package usecase

import (
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
)

type NewsService interface {
	GetNewsSortableColumns() []newsdto.NewsEnumResponse
	GetAllNewsStatuses() []newsdto.NewsEnumResponse
	GetAdminNews(newsID uint) (newsdto.AdminNewsResponse, error)
	GetPublicNews(newsID uint) (newsdto.PublicNewsResponse, error)
	GetAdminNewsList(request newsdto.GetAdminNewsListRequest) ([]newsdto.AdminNewsResponse, int64, error)
	SearchNews(request newsdto.SearchNewsRequest) ([]newsdto.AdminNewsResponse, int64, error)
	GetPublicNewsList(request newsdto.GetPublicNewsListRequest) ([]newsdto.PublicNewsResponse, int64, error)
	CreateNews(request newsdto.CreateNewsRequest) (uint, error)
	EditNews(request newsdto.EditNewsRequest) error
	UpdateNewsStatus(request newsdto.EditNewsStatusRequest) error
	DeleteNewsStatus(request newsdto.DeleteNewsRequest) error
	AddNewsMedia(request newsdto.AddNewsMediaRequest) (uint, error)
	DeleteNewsMedia(request newsdto.AccessMediaRequest) error
	GetNewsMedia(request newsdto.AccessMediaRequest) (string, error)
	LikeNews(request newsdto.GetNewsByCustomer) error
	DislikeNews(request newsdto.GetNewsByCustomer) error
	IsNewsLiked(request newsdto.GetNewsByCustomer) (bool, error)
}
