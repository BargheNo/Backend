package user

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerUserController struct {
	constants   *bootstrap.Constants
	userService service.UserService
}

func NewCustomerUserController(
	constants *bootstrap.Constants,
	userService service.UserService,
) *CustomerUserController {
	return &CustomerUserController{
		constants:   constants,
		userService: userService,
	}
}

func (userController *CustomerUserController) CompleteRegister(ctx *gin.Context) {
	type resetPasswordParams struct {
		Email        string                `form:"email" validate:"omitempty,email"`
		NationalCode string                `form:"nationalCode"`
		ProfilePic   *multipart.FileHeader `form:"profilePic"`
	}
	params := controller.Validated[resetPasswordParams](ctx)
	userID, _ := ctx.Get(userController.constants.Context.ID)

	completeRegisterInfo := userdto.CompleteRegisterRequest{
		UserID:       userID.(uint),
		Email:        params.Email,
		NationalCode: params.NationalCode,
		ProfilePic:   params.ProfilePic,
	}
	userController.userService.CompleteRegister(completeRegisterInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.completeRegister")
	controller.Response(ctx, 200, message, nil)
}

func (userController *CustomerUserController) ResetPassword(ctx *gin.Context) {
	type completeProfileParams struct {
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
	}
	params := controller.Validated[completeProfileParams](ctx)
	userID, _ := ctx.Get(userController.constants.Context.ID)
	resetPasswordInfo := userdto.ResetPasswordRequest{
		ID:       userID.(uint),
		Password: params.Password,
	}
	userController.userService.ResetPassword(resetPasswordInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.resetPassword")
	controller.Response(ctx, 200, message, nil)
}
