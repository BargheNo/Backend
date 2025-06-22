package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralBlogController struct {
	constants   *bootstrap.Constants
	blogService port.BlogService
	pagination  *bootstrap.Pagination
}

func NewGeneralBlogController(
	constants *bootstrap.Constants,
	blogService port.BlogService,
	pagination *bootstrap.Pagination,
) *GeneralBlogController {
	return &GeneralBlogController{
		constants:   constants,
		blogService: blogService,
		pagination:  pagination,
	}
}

func (blogController *GeneralBlogController) GetPosts(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	request := blogdto.GetPublicPostsRequest{
		Offset: offset,
		Limit:  limit,
	}
	posts, err := blogController.blogService.GetGeneralPosts(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", posts)
}

func (blogController *GeneralBlogController) GetCorporationPosts(ctx *gin.Context) {
	type getCorporationPostsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	pagination := controller.GetPagination(ctx, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	params := controller.Validated[getCorporationPostsParams](ctx)

	request := blogdto.GetPublicCorporationPostsRequest{
		CorporationID: params.CorporationID,
		Offset:        offset,
		Limit:         limit,
	}
	posts, err := blogController.blogService.GetCorporationPostsForGeneral(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", posts)
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
