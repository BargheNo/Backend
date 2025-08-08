package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralNewsController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	newsService usecase.NewsService
}

func NewGeneralNewsController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	newsService usecase.NewsService,
) *GeneralNewsController {
	return &GeneralNewsController{
		constants:   constants,
		pagination:  pagination,
		newsService: newsService,
	}
}

func (newsController *GeneralNewsController) GetNewsList(ctx *gin.Context) {
	type getNewsParams struct {
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getNewsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, newsController.pagination.DefaultPage, newsController.pagination.DefaultPageSize)

	getNewsRequest := newsdto.GetPublicNewsListRequest{
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	news, count, err := newsController.newsService.GetPublicNewsList(getNewsRequest)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(news, count, params.Page, params.PageSize)

	controller.Response(ctx, 200, "", data)
}

func (newsController *GeneralNewsController) GetNews(ctx *gin.Context) {
	type getNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[getNewsParams](ctx)

	news, err := newsController.newsService.GetPublicNews(params.NewsID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", news)
}

func (newsController *GeneralNewsController) GetNewsMedia(ctx *gin.Context) {
	type getNewsParams struct {
		NewsID  uint `uri:"newsID" validate:"required"`
		MediaID uint `uri:"mediaID" validate:"required"`
	}
	params := controller.Validated[getNewsParams](ctx)

	mediaParams := newsdto.AccessMediaRequest{
		NewsID:   params.NewsID,
		MediaID:  params.MediaID,
		UserType: enum.UserTypeGuest,
	}
	media, err := newsController.newsService.GetNewsMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", media)
}

func (newsController *GeneralNewsController) GetSortableFields(ctx *gin.Context) {
	columns := newsController.newsService.GetNewsSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
