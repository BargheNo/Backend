package news

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminNewsController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	newsService usecase.NewsService
}

func NewAdminNewsController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	newsService usecase.NewsService,
) *AdminNewsController {
	return &AdminNewsController{
		constants:   constants,
		pagination:  pagination,
		newsService: newsService,
	}
}

func (newsController *AdminNewsController) GetAllNewsStatuses(ctx *gin.Context) {
	statuses := newsController.newsService.GetAllNewsStatuses()
	controller.Response(ctx, 200, "", statuses)
}

func (newsController *AdminNewsController) CreateDraftNews(ctx *gin.Context) {
	type createNewsParams struct {
		Title       string                `json:"title" validate:"required"`
		Content     string                `json:"content"`
		Description string                `json:"description"`
		CoverImage  *multipart.FileHeader `form:"cover_image"`
	}
	params := controller.Validated[createNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	draftNewsParams := newsdto.CreateNewsRequest{
		Title:       params.Title,
		Content:     params.Content,
		Description: params.Description,
		AuthorID:    authorID.(uint),
		Status:      enum.NewsStatusDraft,
		CoverImage:  params.CoverImage,
	}
	news, err := newsController.newsService.CreateNews(draftNewsParams)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createDraftNews")
	controller.Response(ctx, 200, message, news)
}

func (newsController *AdminNewsController) EditNews(ctx *gin.Context) {
	type editNewsParams struct {
		NewsID      uint                  `uri:"newsID" validate:"required"`
		Title       *string               `json:"title"`
		Content     *string               `json:"content"`
		Description *string               `json:"description"`
		CoverImage  *multipart.FileHeader `form:"cover_image"`
		Status      uint                  `json:"status"`
	}
	params := controller.Validated[editNewsParams](ctx)
	authorID, _ := ctx.Get(newsController.constants.Context.ID)

	finalizeNewsParams := newsdto.EditNewsRequest{
		NewsID:      params.NewsID,
		AuthorID:    authorID.(uint),
		Title:       params.Title,
		Content:     params.Content,
		Description: params.Description,
		CoverImage:  params.CoverImage,
		Status:      params.Status,
	}
	if err := newsController.newsService.EditNews(finalizeNewsParams); err != nil {
		panic(err)
	}

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
	if err := newsController.newsService.UpdateNewsStatus(publishParams); err != nil {
		panic(err)
	}

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
	if err := newsController.newsService.UpdateNewsStatus(unpublishParams); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unpublishNews")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *AdminNewsController) GetNewsList(ctx *gin.Context) {
	type getNewsParams struct {
		Status   uint `form:"status" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getNewsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, newsController.pagination.DefaultPage, newsController.pagination.DefaultPageSize)

	getNewsRequest := newsdto.GetAdminNewsListRequest{
		Status: params.Status,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	news, count, err := newsController.newsService.GetAdminNewsList(getNewsRequest)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(news, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (newsController *AdminNewsController) GetNews(ctx *gin.Context) {
	type getNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[getNewsParams](ctx)

	news, err := newsController.newsService.GetAdminNews(params.NewsID)
	if err != nil {
		panic(err)
	}

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
	if err := newsController.newsService.DeleteNewsStatus(deleteParams); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteNews")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *AdminNewsController) AddNewsMedia(ctx *gin.Context) {
	type addMediaParams struct {
		NewsID uint                  `uri:"newsID" validate:"required"`
		Media  *multipart.FileHeader `form:"media" validate:"required"`
	}
	params := controller.Validated[addMediaParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	mediaParams := newsdto.AddNewsMediaRequest{
		NewsID:   params.NewsID,
		AuthorID: userID.(uint),
		Media:    params.Media,
	}
	media, err := newsController.newsService.AddNewsMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addMedia")
	controller.Response(ctx, 200, message, media)
}

func (newsController *AdminNewsController) DeleteNewsMedia(ctx *gin.Context) {
	type deleteMediaParams struct {
		NewsID  uint `uri:"newsID" validate:"required"`
		MediaID uint `uri:"mediaID" validate:"required"`
	}
	params := controller.Validated[deleteMediaParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	mediaParams := newsdto.AccessMediaRequest{
		NewsID:   params.NewsID,
		AuthorID: userID.(uint),
		MediaID:  params.MediaID,
	}
	if err := newsController.newsService.DeleteNewsMedia(mediaParams); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteMedia")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *AdminNewsController) GetNewsMedia(ctx *gin.Context) {
	type getMediaParams struct {
		NewsID  uint `uri:"newsID" validate:"required"`
		MediaID uint `uri:"mediaID" validate:"required"`
	}
	params := controller.Validated[getMediaParams](ctx)

	mediaParams := newsdto.AccessMediaRequest{
		NewsID:   params.NewsID,
		MediaID:  params.MediaID,
		UserType: enum.UserTypeAdmin,
	}
	media, err := newsController.newsService.GetNewsMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", media)
}
