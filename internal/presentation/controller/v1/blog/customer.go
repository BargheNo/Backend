package blog

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerBlogController struct {
	constants   *bootstrap.Constants
	blogService usecase.BlogService
	pagination  *bootstrap.Pagination
}

func NewCustomerBlogController(
	constants *bootstrap.Constants,
	blogService usecase.BlogService,
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

	request := blogdto.GetPostRequest{
		UserID: userID.(uint),
		PostID: params.PostID,
	}
	if err := blogController.blogService.LikePost(request); err != nil {
		panic(err)
	}

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

	request := blogdto.GetPostRequest{
		UserID: userID.(uint),
		PostID: params.PostID,
	}
	if err := blogController.blogService.UnlikePost(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unlikePost")
	controller.Response(ctx, 200, message, nil)
}

func (blogController *CustomerBlogController) IsUserLikedBlog(ctx *gin.Context) {
	type dislikeNewsParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	params := controller.Validated[dislikeNewsParams](ctx)
	userID, _ := ctx.Get(blogController.constants.Context.ID)

	getNewsRequest := blogdto.GetPostRequest{
		PostID: params.PostID,
		UserID: userID.(uint),
	}
	isLiked, err := blogController.blogService.IsBlogLiked(getNewsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", isLiked)
}
