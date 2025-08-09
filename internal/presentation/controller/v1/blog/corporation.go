package blog

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationBlogController struct {
	constants   *bootstrap.Constants
	blogService usecase.BlogService
	pagination  *bootstrap.Pagination
}

func NewCorporationBlogController(
	constants *bootstrap.Constants,
	blogService usecase.BlogService,
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
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		Title         string                `form:"title" validate:"required"`
		Content       string                `form:"content" validate:"required"`
		Description   string                `form:"description" validate:"required"`
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
	blog, err := blogController.blogService.CreatePost(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createPost")
	controller.Response(ctx, 200, message, blog)
}

func (blogController *CorporationBlogController) EditPost(ctx *gin.Context) {
	type editPostParams struct {
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		PostID        uint                  `uri:"postID" validate:"required"`
		Status        uint                  `form:"status"`
		Title         *string               `form:"title"`
		Content       *string               `form:"content"`
		Description   *string               `form:"description"`
		CoverImage    *multipart.FileHeader `form:"cover_image"`
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
	if err := blogController.blogService.EditPost(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.editPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) PublishPost(ctx *gin.Context) {
	type publishPostParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PostID        uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[publishPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.EditPostRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Status:        uint(enum.PostStatusPublished),
		CorporationID: params.CorporationID,
	}
	if err := blogController.blogService.EditPost(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.publishPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) UnpublishPost(ctx *gin.Context) {
	type unpublishPostParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PostID        uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[unpublishPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.EditPostRequest{
		PostID:        params.PostID,
		AuthorID:      authorID.(uint),
		Status:        uint(enum.PostStatusDraft),
		CorporationID: params.CorporationID,
	}
	if err := blogController.blogService.EditPost(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unpublishPost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) DeletePost(ctx *gin.Context) {
	type deletePostParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		PostIDs       []uint `json:"postIDs" validate:"required"`
	}
	params := controller.Validated[deletePostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	deleteParams := blogdto.DeletePostRequest{
		PostIDs:       params.PostIDs,
		AuthorID:      authorID.(uint),
		CorporationID: params.CorporationID,
	}
	if err := blogController.blogService.DeletePost(deleteParams); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deletePost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) AddPostMedia(ctx *gin.Context) {
	type addPostMediaParams struct {
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		PostID        uint                  `uri:"postID" validate:"required"`
		Media         *multipart.FileHeader `form:"media" validate:"required"`
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
		CorporationID uint `uri:"corporationID" validate:"required"`
		PostID        uint `uri:"postID" validate:"required"`
		MediaID       uint `uri:"mediaID" validate:"required"`
	}
	params := controller.Validated[deletePostMediaParams](ctx)
	userID, _ := ctx.Get(blogController.constants.Context.ID)

	mediaParams := blogdto.AccessPostMediaRequest{
		PostID:        params.PostID,
		UserID:        userID.(uint),
		MediaID:       params.MediaID,
		CorporationID: params.CorporationID,
	}
	if err := blogController.blogService.DeletePostMedia(mediaParams); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteMedia")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CorporationBlogController) GetPosts(ctx *gin.Context) {
	type getPostsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		Status        uint `form:"status"`
		Page          int  `form:"page"`
		PageSize      int  `form:"pageSize"`
		SortBy        uint `form:"sortBy"`
		Asc           bool `form:"asc"`
	}
	params := controller.Validated[getPostsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(blogController.constants.Context.ID)

	getPostsRequest := blogdto.GetCorporationPostsRequest{
		UserID:        userID.(uint),
		CorporationID: params.CorporationID,
		Status:        params.Status,
		Offset:        offset,
		Limit:         limit,
		SortBy:        params.SortBy,
		Asc:           params.Asc,
	}
	posts, count, err := blogController.blogService.GetCorporationPosts(getPostsRequest)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(posts, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (blogController *CorporationBlogController) GetPost(ctx *gin.Context) {
	type getPostParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PostID        uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[getPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	getPostRequest := blogdto.GetCorporationPostRequest{
		UserID:        authorID.(uint),
		PostID:        params.PostID,
		CorporationID: params.CorporationID,
	}
	post, err := blogController.blogService.GetCorporationPost(getPostRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", post)
}

func (blogController *CorporationBlogController) GetPostMedia(ctx *gin.Context) {
	type getPostMediaParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PostID        uint `uri:"postID" validate:"required"`
		MediaID       uint `uri:"mediaID" validate:"required"`
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
