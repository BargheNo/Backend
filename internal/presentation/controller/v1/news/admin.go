package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
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
		Status:   enum.NewsStatusDraft,
	}
	news := newsController.newsService.CreateNews(draftNewsParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createDraftNews")
	controller.Response(ctx, 200, message, news)
}

func (newsController *AdminNewsController) EditNews(ctx *gin.Context) {
	type editNewsParams struct {
		NewsID  uint    `uri:"newsID" validate:"required"`
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Status  uint    `json:"status" validate:"required"`
	}
	params := controller.Validated[editNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	finalizeNewsParams := newsdto.EditNewsRequest{
		NewsID:   params.NewsID,
		AuthorID: authorID.(uint),
		Title:    params.Title,
		Content:  params.Content,
		Status:   params.Status,
	}
	newsController.newsService.EditNews(finalizeNewsParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.editNews")
	controller.Response(ctx, 200, message, nil)
}
func (newsController *AdminNewsController) PublishNews(ctx *gin.Context) {
	type publishNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[publishNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	publishParams := newsdto.EditNewsStatusRequest{
		NewsID:   params.NewsID,
		AuthorID: authorID.(uint),
		Status:   uint(enum.NewsStatusActive),
	}
	newsController.newsService.UpdateNewsStatus(publishParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.publishNews")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *AdminNewsController) UnpublishNews(ctx *gin.Context) {
	type unpublishNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[unpublishNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	unpublishParams := newsdto.EditNewsStatusRequest{
		NewsID:   params.NewsID,
		AuthorID: authorID.(uint),
		Status:   uint(enum.NewsStatusDraft),
	}
	newsController.newsService.UpdateNewsStatus(unpublishParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unpublishNews")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *AdminNewsController) GetAllNewsStatuses(ctx *gin.Context) {
	statuses := newsController.newsService.GetAllNewsStatuses()
	controller.Response(ctx, 200, "", statuses)
}

func (newsController *AdminNewsController) GetNewsList(ctx *gin.Context) {
	type getNewsParams struct {
		Statuses []uint `form:"statuses" validate:"required"`
	}
	params := controller.Validated[getNewsParams](ctx)
	pagination := controller.GetPagination(ctx, newsController.pagination.DefaultPage, newsController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	getNewsRequest := newsdto.GetNewsListRequest{
		Statuses: params.Statuses,
		Offset:   offset,
		Limit:    limit,
	}
	news := newsController.newsService.GetNewsList(getNewsRequest)

	controller.Response(ctx, 200, "", news)
}

func (newsController *AdminNewsController) GetNews(ctx *gin.Context) {
	type getNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[getNewsParams](ctx)

	getNewsRequest := newsdto.GetNewsRequest{
		NewsID:   params.NewsID,
		UserType: enum.UserTypeAdmin,
	}
	news := newsController.newsService.GetNews(getNewsRequest)

	controller.Response(ctx, 200, "", news)
}

func (newsController *AdminNewsController) DeleteNews(ctx *gin.Context) {
	type deleteNewsParams struct {
		NewsIDs []uint `uri:"newsIDs" validate:"required"`
	}
	params := controller.Validated[deleteNewsParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	deleteParams := newsdto.DeleteNewsRequest{
		NewsIDs:  params.NewsIDs,
		AuthorID: userID.(uint),
	}
	newsController.newsService.DeleteNewsStatus(deleteParams)

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteNews")
	controller.Response(ctx, 200, message, nil)
}
