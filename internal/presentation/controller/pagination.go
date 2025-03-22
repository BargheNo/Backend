package controller

import (
	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func GetPagination(ctx *gin.Context, defaultPage, defaultPageSize int) PaginationParams {
	params := Validated[PaginationParams](ctx)

	if params.Page <= 0 {
		params.Page = defaultPage
	}

	if params.PageSize <= 0 {
		params.PageSize = defaultPageSize
	}

	return params
}

func (p PaginationParams) GetOffsetLimit() (offset, limit int) {
	return (p.Page - 1) * p.PageSize, p.PageSize
}
