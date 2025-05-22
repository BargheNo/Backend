package blog

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationBlogController struct {
	constants   *bootstrap.Constants
	blogService service.BlogService
}

func NewCorporationBlogController(
	constants *bootstrap.Constants,
	blogService service.BlogService,
) *CorporationBlogController {
	return &CorporationBlogController{
		constants:   constants,
		blogService: blogService,
	}
}

func (blogController *CorporationBlogController) CreatePost(ctx *gin.Context) {
	type createPostParams struct {
		Title         string                `json:"title" validate:"required"`
		Content       string                `json:"content" validate:"required"`
		AuthorID      uint                  `json:"author_id"`
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		CoverImage    *multipart.FileHeader `form:"cover_image"`
	}
	params := controller.Validated[createPostParams](ctx)
	authorID, _ := ctx.Get(blogController.constants.Context.ID)

	request := blogdto.CreatePostRequest{
		Title:         params.Title,
		Content:       params.Content,
		AuthorID:      authorID.(uint),
		CorporationID: params.CorporationID,
		CoverImage:    params.CoverImage,
	}

	blogController.blogService.CreatePost(request)

	trans := controller.GetTranslator(ctx, blogController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createPost")
	controller.Response(ctx, 200, message, nil)
}
