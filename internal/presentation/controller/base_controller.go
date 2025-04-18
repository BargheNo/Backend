package controller

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/localization"
	"github.com/gin-gonic/gin"
)

func GetTranslator(ctx *gin.Context, key string) localization.TranslatorInstance {
	translator, exists := ctx.Get(key)
	if !exists {
		panic("translator not registered!")
	}

	return translator.(localization.TranslatorInstance)
}

func GetLocalizedTemplateFile(ctx *gin.Context, key, persianTemplateFile, englishTemplateFile string) string {
	trans := GetTranslator(ctx, key)
	switch trans.Locale() {
	case "fa_IR":
		return persianTemplateFile
	case "en_US":
		return englishTemplateFile
	default:
		return persianTemplateFile
	}
}

type SortParams struct {
	SortBy string `form:"sortBy"`
	Dir    string `form:"dir"`
}

func GetSort(c *gin.Context, context *bootstrap.Context) SortParams {
	param := Validated[SortParams](c)
	if param.Dir == "" {
		param.Dir = "ASC"
	}
	return param
}
