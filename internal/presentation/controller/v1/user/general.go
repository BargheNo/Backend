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
	jwtService  service.JWTService
}

func NewGeneralUserController(
	constants *bootstrap.Constants,
	userService service.UserService,
	jwtService service.JWTService,
) *GeneralUserController {
	return &GeneralUserController{
		constants:   constants,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (userController *GeneralUserController) BasicRegister(ctx *gin.Context) {
	type registerParams struct {
		FirstName       string `json:"firstName" validate:"required"`
		LastName        string `json:"lastName" validate:"required"`
		Phone           string `json:"phone" validate:"required,e164"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
		AcceptedTerms   bool   `json:"acceptedTerms" validate:"required,eq=true"`
	}
	params := controller.Validated[registerParams](ctx)
	registerInfo := userdto.BasicRegisterRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     params.Phone,
		Password:  params.Password,
	}
	if err := userController.userService.Register(registerInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.userRegister")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) VerifyPhone(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Phone string `json:"phone" validate:"required,e164"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyPhoneParams](ctx)
	verifyOTPInfo := userdto.VerifyPhoneRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}
	if err := userController.userService.VerifyPhone(verifyOTPInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.phoneVerification")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) Login(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Phone    string `json:"phone" validate:"required,e164"`
		Password string `json:"password" validate:"required"`
	}
	params := controller.Validated[verifyPhoneParams](ctx)
	loginInfo := userdto.LoginRequest{
		Phone:    params.Phone,
		Password: params.Password,
	}
	userInfo, err := userController.userService.Login(loginInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.login")
	controller.Response(ctx, 200, message, userInfo)
}

func (userController *GeneralUserController) ForgotPassword(ctx *gin.Context) {
	type forgotPasswordParams struct {
		Phone string `json:"phone" validate:"required,e164"`
	}
	params := controller.Validated[forgotPasswordParams](ctx)
	forgotPasswordInfo := userdto.ForgotPasswordRequest{
		Phone: params.Phone,
	}
	if err := userController.userService.ForgotPassword(forgotPasswordInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.forgotPassword")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) ConfirmOTP(ctx *gin.Context) {
	type verifyOTPParams struct {
		Phone string `json:"phone" validate:"required,e164"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyOTPParams](ctx)
	verifyPhoneInfo := userdto.VerifyPhoneRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}
	userInfo, err := userController.userService.VerifyOTP(verifyPhoneInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.phoneVerification")
	controller.Response(ctx, 200, message, userInfo)
}

func (userController *GeneralUserController) RefreshToken(ctx *gin.Context) {
	type refreshTokenParams struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	params := controller.Validated[refreshTokenParams](ctx)
	claims, err := userController.jwtService.ValidateToken(params.RefreshToken)
	if err != nil {
		panic(err)
	}

	userID := uint(claims["sub"].(float64))
	accessToken, _, err := userController.jwtService.GenerateToken(userID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.refreshToken")
	controller.Response(ctx, 200, message, accessToken)
}
