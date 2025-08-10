package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralBlogController struct {
	constants   *bootstrap.Constants
	blogService usecase.BlogService
	pagination  *bootstrap.Pagination
}

func NewGeneralBlogController(
	constants *bootstrap.Constants,
	blogService usecase.BlogService,
	pagination *bootstrap.Pagination,
) *GeneralBlogController {
	return &GeneralBlogController{
		constants:   constants,
		blogService: blogService,
		pagination:  pagination,
	}
}

func (blogController *GeneralBlogController) GetPosts(ctx *gin.Context) {
	type getPostsParams struct {
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getPostsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)

	request := blogdto.GetPublicPostsRequest{
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	posts, count, err := blogController.blogService.GetGeneralPosts(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(posts, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (blogController *GeneralBlogController) GetCorporationPosts(ctx *gin.Context) {
	type getCorporationPostsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		Page          int  `form:"page"`
		PageSize      int  `form:"pageSize"`
		SortBy        uint `form:"sortBy"`
		Asc           bool `form:"asc"`
	}
	params := controller.Validated[getCorporationPostsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)

	request := blogdto.GetPublicCorporationPostsRequest{
		CorporationID: params.CorporationID,
		Offset:        offset,
		Limit:         limit,
		SortBy:        params.SortBy,
		Asc:           params.Asc,
	}
	posts, count, err := blogController.blogService.GetCorporationPostsForGeneral(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(posts, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (blogController *GeneralBlogController) GetPost(ctx *gin.Context) {
	type getPostParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[getPostParams](ctx)

	post, err := blogController.blogService.GetGeneralPost(params.PostID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", post)
}

func (blogController *GeneralBlogController) GetPostMedia(ctx *gin.Context) {
	type getPostMediaParams struct {
		PostID  uint `uri:"postID" validate:"required"`
		MediaID uint `uri:"mediaID" validate:"required"`
	}
	params := controller.Validated[getPostMediaParams](ctx)

	mediaParams := blogdto.AccessPostMediaRequest{
		PostID:   params.PostID,
		MediaID:  params.MediaID,
		UserType: enum.UserTypeGuest,
	}
	media, err := blogController.blogService.GetPostMedia(mediaParams)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", media)
}

func (blogController *GeneralBlogController) GetSortableFields(ctx *gin.Context) {
	columns := blogController.blogService.GetBlogSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
