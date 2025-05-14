package service

import (
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
)

type NewsService interface {
	GetAllNewsStatuses() []newsdto.NewsStatusesResponse
	GetNews(newsID uint) newsdto.NewsResponse
	GetNewsList(request newsdto.GetNewsListRequest) []newsdto.NewsResponse
	CreateNews(request newsdto.CreateNewsRequest) newsdto.NewsResponse
	EditNews(request newsdto.EditNewsRequest)
	UpdateNewsStatus(request newsdto.EditNewsStatusRequest)
	DeleteNewsStatus(request newsdto.DeleteNewsRequest)
}
