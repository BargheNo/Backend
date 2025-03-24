package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	loggerimpl "github.com/BargheNo/Backend/internal/application/adapter/logger"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

const (
	genericError = "errors.generic"
)

type RecoveryMiddleware struct {
	constants *bootstrap.Constants
}

func NewRecovery(constants *bootstrap.Constants) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		constants: constants,
	}
}

func (recovery RecoveryMiddleware) Recovery(ctx *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				recovery.handleRecoveredError(ctx, err)
				ctx.Abort()
			}
		}
	}()

	ctx.Next()
}

func (recovery RecoveryMiddleware) handleRecoveredError(ctx *gin.Context, err error) {
	if validationErrors, ok := err.(exception.ValidationErrors); ok {
		handleValidationError(ctx, validationErrors, recovery.constants.Context.Translator)
	} else if bindingError, ok := err.(exception.BindingError); ok {
		handleBindingError(ctx, bindingError, recovery.constants.Context.Translator)
	} else if rateLimitError, ok := err.(*exception.RateLimitError); ok {
		handleRateLimitError(ctx, *rateLimitError, recovery.constants.Context.Translator)
	} else if conflictErrors, ok := err.(exception.ConflictErrors); ok {
		handleConflictError(ctx, conflictErrors, recovery.constants.Context.Translator)
	} else if authError, ok := err.(*exception.AuthError); ok {
		handleAuthError(ctx, *authError, recovery.constants.Context.Translator)
	} else if notFoundError, ok := err.(exception.NotFoundError); ok {
		handleNotFoundError(ctx, notFoundError, recovery.constants.Context.Translator)
	} else if forbiddenError, ok := err.(exception.ForbiddenError); ok {
		handleForbiddenError(ctx, forbiddenError, recovery.constants.Context.Translator)
	} else {
		unhandledErrors(ctx, err, recovery.constants.Context.Translator)
	}
}

func handleValidationError(ctx *gin.Context, validationErrors exception.ValidationErrors, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, validationError := range validationErrors.Errors {
		if _, ok := errorMessages[validationError.Field]; !ok {
			errorMessages[validationError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(validationError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", validationError.Tag), fieldName)
		errorMessages[validationError.Field][validationError.Tag] = message
	}

	controller.Response(ctx, 422, errorMessages, nil)
}

func handleBindingError(ctx *gin.Context, bindingError exception.BindingError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	message, _ := trans.Translate(genericError)

	if numError, ok := bindingError.Err.(*strconv.NumError); ok {
		message, _ = trans.Translate("errors.numeric", numError.Num)
	} else if bindingError == http.ErrMissingFile {
		message, _ = trans.Translate("errors.fileRequired")
	}

	controller.Response(ctx, 400, message, nil)
}

func handleRateLimitError(ctx *gin.Context, rateLimitError exception.RateLimitError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)

	message, _ := trans.Translate(genericError)
	switch rateLimitError.Type {
	case exception.ErrorTypeRequestRateLimit:
		message, _ = trans.Translate("errors.rateLimit")
	case exception.ErrorTypeConcurrentInstallLimit:
		message, _ = trans.Translate("errors.installRateLimit")
	}

	controller.Response(ctx, 429, message, nil)
}

func handleConflictError(ctx *gin.Context, conflictErrors exception.ConflictErrors, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, conflictError := range conflictErrors.Errors {
		if _, ok := errorMessages[conflictError.Field]; !ok {
			errorMessages[conflictError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(conflictError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", conflictError.Tag), fieldName)
		errorMessages[conflictError.Field][conflictError.Tag] = message
	}

	controller.Response(ctx, 409, errorMessages, nil)
}

func handleAuthError(ctx *gin.Context, authError exception.AuthError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)

	message, _ := trans.Translate(genericError)
	switch authError.Type {
	case exception.ErrorTypeInvalidCredentials:
		message, _ = trans.Translate("errors.invalidAuthCredentials")
	case exception.ErrorTypeExpiredToken:
		message, _ = trans.Translate("errors.expiredAuthToken")
	case exception.ErrorTypeInvalidToken:
		message, _ = trans.Translate("errors.invalidAuthToken")
	case exception.ErrorTypeUnauthorized:
		message, _ = trans.Translate("errors.unauthorized")
	}

	controller.Response(ctx, 401, message, nil)
}

func handleNotFoundError(ctx *gin.Context, notFoundError exception.NotFoundError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	itemName, _ := trans.Translate(notFoundError.Item)
	message, _ := trans.Translate("errors.notFound", itemName)
	controller.Response(ctx, 404, message, nil)
}

func handleForbiddenError(ctx *gin.Context, forbiddenError exception.ForbiddenError, transKey string) {
	// maybe add without field later -> add type and switch case between them
	trans := controller.GetTranslator(ctx, transKey)
	ResourceName, _ := trans.Translate(forbiddenError.Resource)
	message, _ := trans.Translate("errors.forbiddenError", ResourceName)
	controller.Response(ctx, 403, message, nil)
}

func unhandledErrors(ctx *gin.Context, err error, transKey string) {
	loggerimpl.GetLogger().Error("unhandled error recovery middleware", logger.Error("error:", err))
	trans := controller.GetTranslator(ctx, transKey)
	errorMessage, _ := trans.Translate(genericError)

	controller.Response(ctx, 500, errorMessage, nil)
}
