package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
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

func (userController *GeneralUserController) BasicRegister(ctx *gin.Context) {
	type registerParams struct {
		FirstName       string `json:"firstName" validate:"required"`
		LastName        string `json:"lastName" validate:"required"`
		Phone           string `json:"phone" validate:"required,e164"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
		IsAcceptTerms   bool   `json:"isAcceptTerms" validate:"required,eq=true"`
	}
	params := controller.Validated[registerParams](ctx, &userController.constants.Context)
	registerInfo := userdto.BasicRegisterRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     params.Phone,
		Password:  params.Password,
	}
	userController.userService.Register(registerInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.userRegister")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) VerifyPhone(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Phone string `json:"phone" validate:"required,e164"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyPhoneParams](ctx, &userController.constants.Context)
	verifyOTPInfo := userdto.VerifyPhoneRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}
	userController.userService.VerifyPhone(verifyOTPInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.phoneVerification")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) VerifyEmail(ctx *gin.Context) {
	// some code here ...
}

func (userController *GeneralUserController) Login(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Phone    string `json:"phone" validate:"required,e164"`
		Password string `json:"password" validate:"required"`
	}
	params := controller.Validated[verifyPhoneParams](ctx, &userController.constants.Context)
	verifyOTPInfo := userdto.LoginRequest{
		Phone:    params.Phone,
		Password: params.Password,
	}
	userInfo := userController.userService.Login(verifyOTPInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.login")
	controller.Response(ctx, 200, message, userInfo)
}

func (userController *GeneralUserController) ForgotPassword(ctx *gin.Context) {
	type forgotPasswordParams struct {
		Phone string `json:"phone" validate:"required,e164"`
	}
	params := controller.Validated[forgotPasswordParams](ctx, &userController.constants.Context)
	forgotPasswordInfo := userdto.ForgotPassword{
		Phone: params.Phone,
	}
	userController.userService.ForgotPassword(forgotPasswordInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.forgotPassword")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) ConfirmOTP(ctx *gin.Context) {
	type verifyOTPParams struct {
		Phone string `json:"phone" validate:"required,e164"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyOTPParams](ctx, &userController.constants.Context)
	verifyPhoneInfo := userdto.VerifyPhoneRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}
	userInfo := userController.userService.VerifyOTP(verifyPhoneInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.phoneVerification")
	controller.Response(ctx, 200, message, userInfo)
}

func (userController *GeneralUserController) ResetPassword(ctx *gin.Context) {
	// some code here ...
}

func (userController *GeneralUserController) RefreshToken(ctx *gin.Context) {
	// some code here ...
}
