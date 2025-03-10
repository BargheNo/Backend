package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	constants   *bootstrap.Constants
	userService service.UserService
}

func NewGeneralUserController(
	constants *bootstrap.Constants,
	userService service.UserService,
) *GeneralUserController {
	return &GeneralUserController{
		constants:   constants,
		userService: userService,
	}
}

func (controller *GeneralUserController) BasicRegister(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) RegisterPhone(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) RegisterEmail(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) Login(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) ForgotPassword(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) ResetPassword(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) ConfirmOTP(ctx *gin.Context) {
	// some code here ...
}

func (controller *GeneralUserController) RefreshToken(ctx *gin.Context) {
	// some code here ...
}
