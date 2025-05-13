package service

import (
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
)

type NewsService interface {
	CreateNews(request newsdto.CreateNewsRequest) newsdto.NewsResponse
	EditNews(request newsdto.EditNewsRequest)
	UpdateNewsStatus(request newsdto.EditNewsStatusRequest)
}
