package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
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
	type createNewsParams struct {
		Title   string `json:"title" validate:"required"`
		Content string `json:"content"`
	}
	params := controller.Validated[createNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	draftNewsParams := newsdto.CreateNewsRequest{
		Title:    params.Title,
		Content:  params.Content,
		AuthorID: authorID.(uint),
	}
	newsController.newsService.CreateNews(draftNewsParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createDraftNews")
	controller.Response(ctx, 200, message, nil)
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
