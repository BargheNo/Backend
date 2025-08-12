package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService usecase.CorporationService
}

func NewAdminCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService usecase.CorporationService,
) *AdminCorporationController {
	return &AdminCorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
	}
}

func (corporationController *AdminCorporationController) GetCorporationStatus(ctx *gin.Context) {
	statuses := corporationController.corporationService.GetCorporationStatuses()
	controller.Response(ctx, 200, "", statuses)
}

func (corporationController *AdminCorporationController) GetCorporations(ctx *gin.Context) {
	type getCorporationsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)

	listInfo := corporationdto.GetCorporationsByAdminRequest{
		Status: params.Status,
		Limit:  limit,
		Offset: offset,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	corporations, count, err := corporationController.corporationService.GetCorporationsByAdmin(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(corporations, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (corporationController *AdminCorporationController) SearchCorporations(ctx *gin.Context) {
	type searchCorporationsParams struct {
		Query    string `form:"query" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}

	params := controller.Validated[searchCorporationsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)

	request := corporationdto.SearchCorporationsRequest{
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	corporations, count, err := corporationController.corporationService.SearchCorporations(request)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(corporations, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (corporationController *AdminCorporationController) GetCorporation(ctx *gin.Context) {
	type getCorporationsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	corporation, err := corporationController.corporationService.GetCorporationByAdmin(params.CorporationID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", corporation)
}

func (corporationController *AdminCorporationController) GetReviewActions(ctx *gin.Context) {
	actions := corporationController.corporationService.GetReviewActions()
	controller.Response(ctx, 200, "", actions)
}

func (corporationController *AdminCorporationController) GetCorporationReviews(ctx *gin.Context) {
	type getCorporationsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	reviews, err := corporationController.corporationService.GetCorporationReviewsByAdmin(params.CorporationID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", reviews)
}

func (corporationController *AdminCorporationController) ApproveCorporationRequest(ctx *gin.Context) {
	type approveCorporationParams struct {
		CorporationID uint    `uri:"corporationID" validate:"required"`
		Reason        *string `json:"reason"`
		Notes         *string `json:"notes"`
	}
	params := controller.Validated[approveCorporationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	request := corporationdto.HandleCorporationActionRequest{
		CorporationID: params.CorporationID,
		ReviewerID:    userID.(uint),
		ActionID:      uint(enum.ReviewActionApproved),
		Reason:        params.Reason,
		Notes:         params.Notes,
	}
	if err := corporationController.corporationService.ApproveCorporationRegistration(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.approveCorporation")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *AdminCorporationController) RejectCorporationRequest(ctx *gin.Context) {
	type approveCorporationParams struct {
		CorporationID uint    `uri:"corporationID" validate:"required"`
		ActionID      uint    `json:"action" validate:"required"`
		Reason        *string `json:"reason"`
		Notes         *string `json:"notes"`
	}
	params := controller.Validated[approveCorporationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	request := corporationdto.HandleCorporationActionRequest{
		CorporationID: params.CorporationID,
		ReviewerID:    userID.(uint),
		ActionID:      params.ActionID,
		Reason:        params.Reason,
		Notes:         params.Notes,
	}
	if err := corporationController.corporationService.RejectCorporationRegistration(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectCorporation")
	controller.Response(ctx, 200, message, nil)
}
