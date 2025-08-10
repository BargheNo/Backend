package news

import (
	"github.com/BargheNo/Backend/bootstrap"
	newsdto "github.com/BargheNo/Backend/internal/application/dto/news"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerNewsController struct {
	constants   *bootstrap.Constants
	newsService usecase.NewsService
}

func NewCustomerNewsController(
	constants *bootstrap.Constants,
	newsService usecase.NewsService,
) *CustomerNewsController {
	return &CustomerNewsController{
		constants:   constants,
		newsService: newsService,
	}
}

func (newsController *CustomerNewsController) LikeNews(ctx *gin.Context) {
	type likeNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[likeNewsParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	getNewsRequest := newsdto.GetNewsByCustomer{
		NewsID: params.NewsID,
		UserID: userID.(uint),
	}
	if err := newsController.newsService.LikeNews(getNewsRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.likePost")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *CustomerNewsController) DislikeNews(ctx *gin.Context) {
	type dislikeNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[dislikeNewsParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	getNewsRequest := newsdto.GetNewsByCustomer{
		NewsID: params.NewsID,
		UserID: userID.(uint),
	}
	if err := newsController.newsService.DislikeNews(getNewsRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, newsController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unlikePost")
	controller.Response(ctx, 200, message, nil)
}

func (newsController *CustomerNewsController) IsUserLikedNews(ctx *gin.Context) {
	type dislikeNewsParams struct {
		NewsID uint `uri:"newsID" validate:"required"`
	}
	params := controller.Validated[dislikeNewsParams](ctx)
	userID, _ := ctx.Get(newsController.constants.Context.ID)

	getNewsRequest := newsdto.GetNewsByCustomer{
		NewsID: params.NewsID,
		UserID: userID.(uint),
	}
	isLiked, err := newsController.newsService.IsNewsLiked(getNewsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", isLiked)
}
