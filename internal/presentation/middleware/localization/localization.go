package localization

import (
	"net/http"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/localization"
	"github.com/gin-gonic/gin"
)

type LocalizationMiddleware struct {
	constants  *bootstrap.Constants
	translator localization.Translator
}

func NewLocalization(constants *bootstrap.Constants, translator localization.Translator) *LocalizationMiddleware {
	return &LocalizationMiddleware{
		translator: translator,
		constants:  constants,
	}
}

func (lm LocalizationMiddleware) Localization(c *gin.Context) {
	locale := getLocale(c.Request)

	translatorInstance := lm.translator.GetTranslator(locale)
	c.Set(lm.constants.Context.Translator, translatorInstance)

	c.Next()
}

func getLocale(request *http.Request) string {
	return request.Header.Get("Accept-Language")
}
