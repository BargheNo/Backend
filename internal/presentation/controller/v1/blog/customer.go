package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerBlogController struct {
	constants   *bootstrap.Constants
	blogService port.BlogService
	pagination  *bootstrap.Pagination
}

func NewCustomerBlogController(
	constants *bootstrap.Constants,
	blogService port.BlogService,
	pagination *bootstrap.Pagination,
) *CustomerBlogController {
	return &CustomerBlogController{
		constants:   constants,
		blogService: blogService,
		pagination:  pagination,
	}
}

func (blogController *CustomerBlogController) LikePost(ctx *gin.Context) {
	type likePostParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[likePostParams](ctx)

	userID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.LikePostRequest{
		UserID: userID.(uint),
		PostID: params.PostID,
	}

	blogController.blogService.LikePost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.likePost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CustomerBlogController) UnlikePost(ctx *gin.Context) {
	type unlikePostParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[unlikePostParams](ctx)

	userID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.LikePostRequest{
		UserID: userID.(uint),
		PostID: params.PostID,
	}

	blogController.blogService.UnlikePost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unlikePost")
	controller.Response(ctx, 200, message, nil)
}
