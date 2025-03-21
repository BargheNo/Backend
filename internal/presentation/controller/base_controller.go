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

type PaginationParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func GetPagination(c *gin.Context, context *bootstrap.Context) PaginationParams {
	param := Validated[PaginationParams](c, context)

	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	return param
}

type SortParams struct {
	SortBy string `form:"sortBy"`
	Dir    string `form:"dir"`
}

func GetSort(c *gin.Context, context *bootstrap.Context) SortParams {
	param := Validated[SortParams](c, context)
	if param.Dir == "" {
		param.Dir = "ASC"
	}
	return param
}
