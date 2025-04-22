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

func (userController *CustomerUserController) GetMyProfile(ctx *gin.Context) {
	userID, _ := ctx.Get(userController.constants.Context.ID)
	profile := userController.userService.GetUserCredential(userID.(uint))
	controller.Response(ctx, 200, "", profile)
}

func (userController *CustomerUserController) CompleteRegister(ctx *gin.Context) {
	type resetPasswordParams struct {
		Email        string                `form:"email" validate:"omitempty,email"`
		NationalCode string                `form:"nationalCode"`
		ProfilePic   *multipart.FileHeader `form:"profilePic"`
	}
	params := controller.Validated[resetPasswordParams](ctx)
	userID, _ := ctx.Get(userController.constants.Context.ID)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)

	templateFile := controller.GetLocalizedTemplateFile(ctx, userController.constants.Context.Translator, userController.constants.EmailTemplates.PersianFileName, userController.constants.EmailTemplates.EnglishFileName)
	emailSubject, _ := trans.Translate("emailSubject.emailConfirmation")
	completeRegisterInfo := userdto.CompleteRegisterRequest{
		UserID:       userID.(uint),
		Email:        params.Email,
		NationalCode: params.NationalCode,
		ProfilePic:   params.ProfilePic,
		TemplateFile: "email_confirmation/" + templateFile,
		EmailSubject: emailSubject,
	}
	userController.userService.CompleteRegister(completeRegisterInfo)

	message, _ := trans.Translate("successMessage.completeRegister")
	controller.Response(ctx, 200, message, nil)
}

func (userController *CustomerUserController) VerifyEmail(ctx *gin.Context) {
	type verifyEmailParams struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyEmailParams](ctx)
	userID, _ := ctx.Get(userController.constants.Context.ID)

	verifyOTPInfo := userdto.VerifyEmailRequest{
		UserID: userID.(uint),
		Email:  params.Email,
		OTP:    params.OTP,
	}
	userController.userService.VerifyEmail(verifyOTPInfo)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.emailVerification")
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

func (userController *CustomerUserController) UpdateProfile(ctx *gin.Context) {
	type updateProfileParams struct {
		FirstName    *string               `form:"firstName" validate:"omitempty"`
		LastName     *string               `form:"lastName" validate:"omitempty"`
		Email        *string               `form:"email" validate:"omitempty,email"`
		NationalCode *string               `form:"nationalCode" validate:"omitempty"`
		ProfilePic   *multipart.FileHeader `form:"profilePic" validate:"omitempty"`
	}
	params := controller.Validated[updateProfileParams](ctx)
	userID, _ := ctx.Get(userController.constants.Context.ID)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)

	templateFile := controller.GetLocalizedTemplateFile(ctx, userController.constants.Context.Translator, userController.constants.EmailTemplates.PersianFileName, userController.constants.EmailTemplates.EnglishFileName)
	emailSubject, _ := trans.Translate("emailSubject.emailConfirmation")
	profileInfo := userdto.UpdateProfileRequest{
		UserID:       userID.(uint),
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		NationalCode: params.NationalCode,
		ProfilePic:   params.ProfilePic,
		TemplateFile: "email_confirmation/" + templateFile,
		EmailSubject: emailSubject,
	}
	userController.userService.UpdateProfile(profileInfo)

	message, _ := trans.Translate("successMessage.updateProfile")
	controller.Response(ctx, 200, message, nil)
}
