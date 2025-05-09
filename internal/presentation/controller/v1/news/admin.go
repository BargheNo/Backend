package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gin-gonic/gin"
)

type AdminNewsController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	newsService service.NewsService
}

func NewAdminNewsController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	newsService service.NewsService,
) *AdminNewsController {
	return &AdminNewsController{
		constants:   constants,
		pagination:  pagination,
		newsService: newsService,
	}
}

func (newsController *AdminNewsController) CreateDraftNews(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *AdminNewsController) CreateFinalizeNews(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *AdminNewsController) EditNews(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *AdminNewsController) GetNewsList(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *AdminNewsController) GetNews(ctx *gin.Context) {
	// some codes here ...
}

func (newsController *AdminNewsController) DeleteNews(ctx *gin.Context) {
	// some codes here ...
}
