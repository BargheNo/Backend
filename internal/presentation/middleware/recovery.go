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

type RecoveryMiddleware struct {
	constants *bootstrap.Constants
}

func NewRecovery(constants *bootstrap.Constants) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		constants: constants,
	}
}

func (recovery RecoveryMiddleware) Recovery(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				recovery.handleRecoveredError(c, err)
				c.Abort()
			}
		}
	}()

	c.Next()
}

func (recovery RecoveryMiddleware) handleRecoveredError(c *gin.Context, err error) {
	if validationErrors, ok := err.(exception.ValidationErrors); ok {
		handleValidationError(c, validationErrors, recovery.constants.Context.Translator)
	} else if bindingError, ok := err.(exception.BindingError); ok {
		handleBindingError(c, bindingError, recovery.constants.Context.Translator)
	} else if _, ok := err.(exception.RateLimitError); ok {
		handleRateLimitError(c, recovery.constants.Context.Translator)
	} else if conflictErrors, ok := err.(exception.ConflictErrors); ok {
		handleConflictError(c, conflictErrors, recovery.constants.Context.Translator)
	} else {
		unhandledErrors(c, err, recovery.constants.Context.Translator)
	}
}

func handleValidationError(c *gin.Context, validationErrors exception.ValidationErrors, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, validationError := range validationErrors.Errors {
		if _, ok := errorMessages[validationError.Field]; !ok {
			errorMessages[validationError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(validationError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", validationError.Tag), fieldName)
		errorMessages[validationError.Field][validationError.Tag] = message
	}

	controller.Response(c, 422, errorMessages, nil)
}

func handleBindingError(c *gin.Context, bindingError exception.BindingError, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	message, _ := trans.Translate("errors.generic")

	if numError, ok := bindingError.Err.(*strconv.NumError); ok {
		message, _ = trans.Translate("errors.numeric", numError.Num)
	} else if bindingError == http.ErrMissingFile {
		message, _ = trans.Translate("errors.fileRequired")
	}

	controller.Response(c, 400, message, nil)
}

func handleRateLimitError(c *gin.Context, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	message, _ := trans.Translate("errors.rateLimitExceed")
	controller.Response(c, 429, message, nil)
}

func handleConflictError(c *gin.Context, conflictErrors exception.ConflictErrors, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, conflictError := range conflictErrors.Errors {
		if _, ok := errorMessages[conflictError.Field]; !ok {
			errorMessages[conflictError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(conflictError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", conflictError.Tag), fieldName)
		errorMessages[conflictError.Field][conflictError.Tag] = message
	}

	controller.Response(c, 422, errorMessages, nil)
}

func unhandledErrors(c *gin.Context, err error, transKey string) {
	loggerimpl.GetLogger().Error("unhandled error recovery middleware", logger.Error("error:", err))
	trans := controller.GetTranslator(c, transKey)
	errorMessage, _ := trans.Translate("errors.generic")

	controller.Response(c, 500, errorMessage, nil)
}
