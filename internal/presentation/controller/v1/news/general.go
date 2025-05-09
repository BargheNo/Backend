package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gin-gonic/gin"
)

type GeneralNewsController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	newsService service.NewsService
}

func NewGeneralNewsController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	newsService service.NewsService,
) *GeneralNewsController {
	return &GeneralNewsController{
		constants:   constants,
		pagination:  pagination,
		newsService: newsService,
	}
}

func (newsController *GeneralNewsController) GetNewsList(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *GeneralNewsController) GetNews(ctx *gin.Context) {
	// some codes here ...
}
