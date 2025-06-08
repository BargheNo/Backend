package blog

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationBlogController struct {
	constants   *bootstrap.Constants
	blogService service.BlogService
	pagination  *bootstrap.Pagination
}

func NewCorporationBlogController(
	constants *bootstrap.Constants,
	blogService service.BlogService,
	pagination *bootstrap.Pagination,
) *CorporationBlogController {
	return &CorporationBlogController{
		constants:   constants,
		blogService: blogService,
		pagination:  pagination,
	}
}

func (blogController *CorporationBlogController) CreateDraftPost(ctx *gin.Context) {
	type createPostParams struct {
		Title         string                `json:"title" validate:"required"`
		Content       string                `json:"content" validate:"required"`
		Description   string                `json:"description" validate:"required"`
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		CoverImage    *multipart.FileHeader `form:"cover_image"`
	}
	params := controller.Validated[createPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.CreatePostRequest{
		Title:         params.Title,
		Content:       params.Content,
		Description:   params.Description,
		AuthorID:      authorID.(uint),
		CorporationID: params.CorporationID,
		CoverImage:    params.CoverImage,
		Status:        enum.PostStatusDraft,
	}

	blogController.blogService.CreatePost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) EditPost(ctx *gin.Context) {
	type editPostParams struct {
		PostID        uint                  `uri:"postID" validate:"required"`
		Title         *string               `json:"title"`
		Content       *string               `json:"content"`
		Description   *string               `json:"description"`
		CoverImage    *multipart.FileHeader `form:"cover_image"`
		Status        uint                  `json:"status"`
		CorporationID uint                  `uri:"corporationID" validate:"required"`
	}

	params := controller.Validated[editPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.EditPostRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Title:         params.Title,
		Content:       params.Content,
		Description:   params.Description,
		CoverImage:    params.CoverImage,
		Status:        params.Status,
		CorporationID: params.CorporationID,
	}

	blogController.blogService.EditPost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.editPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) PublishPost(ctx *gin.Context) {
	type publishPostParams struct {
		PostID        uint `uri:"postID" validate:"required"`
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[publishPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.EditPostRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Status:        uint(enum.PostStatusPublished),
		CorporationID: params.CorporationID,
	}
	blogController.blogService.EditPost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.publishPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) UnpublishPost(ctx *gin.Context) {
	type unpublishPostParams struct {
		PostID        uint `uri:"postID" validate:"required"`
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[unpublishPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.EditPostRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Status:        uint(enum.PostStatusDraft),
		CorporationID: params.CorporationID,
	}
	blogController.blogService.EditPost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unpublishPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) DeletePost(ctx *gin.Context) {
	type deletePostParams struct {
		PostIDs       []uint `json:"postIDs" validate:"required"`
		CorporationID uint   `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[deletePostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	deletParams := blogdto.DeletePostRequest{
		PostIDs:       params.PostIDs,
		AuthorID:      authorID.(uint),
		CorporationID: params.CorporationID,
	}
	blogController.blogService.DeletePost(deletParams)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deletePost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) AddPostMedia(ctx *gin.Context) {
	type addPostMediaParams struct {
		PostID        uint                  `uri:"postID" validate:"required"`
		Media         *multipart.FileHeader `form:"media" validate:"required"`
		CorporationID uint                  `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[addPostMediaParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	mediaParams := blogdto.AddPostMediaRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Media:         params.Media,
		CorporationID: params.CorporationID,
	}

	mediaID, err := blogController.blogService.AddPostMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addMedia")
	controller.Response(ctx, 200, message, mediaID)
}

func (blogController *CorporationBlogController) DeletePostMedia(ctx *gin.Context) {
	type deletePostMediaParams struct {
		PostID        uint `uri:"postID" validate:"required"`
		MediaID       uint `uri:"mediaID" validate:"required"`
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[deletePostMediaParams](ctx)
	userID, _ := ctx.Get(blogController.constants.Context.ID)

	mediaParams := blogdto.AccessPostMediaRequest{
		PostID:        params.PostID,
		UserID:        userID.(uint),
		MediaID:       params.MediaID,
		CorporationID: params.CorporationID,
	}
	blogController.blogService.DeletePostMedia(mediaParams)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteMedia")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) GetPosts(ctx *gin.Context) {
	type getPostsParams struct {
		Statuses      []uint `form:"statuses" validate:"required"`
		CorporationID uint   `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getPostsParams](ctx)
	pagination := controller.GetPagination(ctx, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	userID, _ := ctx.Get(blogController.constants.Context.ID)

	getPostsRequest := blogdto.GetPostsRequest{
		CorporationID: params.CorporationID,
		Statuses:      params.Statuses,
		Offset:        offset,
		Limit:         limit,
		UserID:        userID.(uint),
	}
	posts, err := blogController.blogService.GetCorporationPosts(getPostsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", posts)
}

func (blogController *CorporationBlogController) GetPost(ctx *gin.Context) {
	type getPostParams struct {
		PostID        uint `uri:"postID" validate:"required"`
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	getPostRequest := blogdto.GetPostRequest{
		UserID:        authorID.(uint),
		PostID:        params.PostID,
		CorporationID: params.CorporationID,
		UserType:      enum.UserTypeCorporation,
	}
	post, err := blogController.blogService.GetCorporationPost(getPostRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", post)
}

func (blogController *CorporationBlogController) GetPostMedia(ctx *gin.Context) {
	type getPostMediaParams struct {
		PostID        uint `uri:"postID" validate:"required"`
		MediaID       uint `uri:"mediaID" validate:"required"`
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getPostMediaParams](ctx)
	userID, _ := ctx.Get(blogController.constants.Context.ID)

	mediaParams := blogdto.AccessPostMediaRequest{
		PostID:        params.PostID,
		UserID:        userID.(uint),
		MediaID:       params.MediaID,
		CorporationID: params.CorporationID,
	}

	media, err := blogController.blogService.GetPostMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", media)
}
