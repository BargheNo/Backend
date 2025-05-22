package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralBlogController struct {
	constants   *bootstrap.Constants
	blogService service.BlogService
}

func NewGeneralBlogController(
	constants *bootstrap.Constants,
	blogService service.BlogService,
) *GeneralBlogController {
	return &GeneralBlogController{
		constants:   constants,
		blogService: blogService,
	}
}

func (blogController *GeneralBlogController) GetCorporationPosts(ctx *gin.Context) {
	type getCorporationPostsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationPostsParams](ctx)

	posts, err := blogController.blogService.GetCorporationPosts(params.CorporationID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", posts)
}
