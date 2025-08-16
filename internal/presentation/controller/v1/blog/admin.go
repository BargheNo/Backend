package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBlogController struct {
	constants   *bootstrap.Constants
	blogService usecase.BlogService
}

func NewAdminBlogController(
	constants *bootstrap.Constants,
	blogService usecase.BlogService,
) *AdminBlogController {
	return &AdminBlogController{
		constants:   constants,
		blogService: blogService,
	}
}

func (blogController *AdminBlogController) DeletePost(ctx *gin.Context) {
	type deletePostParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[deletePostParams](ctx)

	blogController.blogService.DeletePostByAdmin(params.PostID)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deletePost")
	controller.Response(ctx, 200, message, nil)
}
