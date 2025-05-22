package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralBlogController struct {
	constants   *bootstrap.Constants
	blogService service.BlogService
	pagination  *bootstrap.Pagination
}

func NewGeneralBlogController(
	constants *bootstrap.Constants,
	blogService service.BlogService,
	pagination *bootstrap.Pagination,
) *GeneralBlogController {
	return &GeneralBlogController{
		constants:   constants,
		blogService: blogService,
		pagination:  pagination,
	}
}

func (blogController *GeneralBlogController) GetCorporationPosts(ctx *gin.Context) {
	type getCorporationPostsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	pagination := controller.GetPagination(ctx, blogController.pagination.DefaultPage, blogController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	params := controller.Validated[getCorporationPostsParams](ctx)

	request := blogdto.GetCorporationPostsRequest{
		CorporationID: params.CorporationID,
		Offset:        offset,
		Limit:         limit,
	}
	posts, err := blogController.blogService.GetCorporationPosts(request)
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

	post := blogController.blogService.GetPost(params.PostID)

	controller.Response(ctx, 200, "", post)
}
