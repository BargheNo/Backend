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
		Query    string `form:"query"`
		Status   uint   `form:"status"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)

	listInfo := corporationdto.GetCorporationsByAdminRequest{
		Status: params.Status,
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
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

func (corporationController *AdminCorporationController) GetStaffList(ctx *gin.Context) {
	type getStaffParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Query         string `form:"query"`
		Status        uint   `form:"status"`
		Page          int    `form:"page"`
		PageSize      int    `form:"pageSize"`
		SortBy        uint   `form:"sortBy"`
		Asc           bool   `form:"asc"`
	}
	params := controller.Validated[getStaffParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)

	request := corporationdto.GetStaffList{
		CorporationID: params.CorporationID,
		Query:         params.Query,
		Status:        params.Status,
		Offset:        offset,
		Limit:         limit,
		SortBy:        params.SortBy,
		Asc:           params.Asc,
	}
	staffs, count, err := corporationController.corporationService.GetStaffList(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(staffs, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (corporationController *AdminCorporationController) CreateCorporationStaff(ctx *gin.Context) {
	type addStaffParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Phone         string `json:"phone" validate:"required,e164"`
		RoleIDs       []uint `json:"roleIDs" validate:"required"`
	}
	params := controller.Validated[addStaffParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	request := corporationdto.AddStaffRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		StaffPhone:    params.Phone,
		RoleIDs:       params.RoleIDs,
	}
	if err := corporationController.corporationService.AddStaff(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addStaff")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *AdminCorporationController) EditCorporationStaff(ctx *gin.Context) {
	type editStaffParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		StaffID       uint   `uri:"staffID" validate:"required"`
		Status        *uint  `json:"status"`
		RoleIDs       []uint `json:"roleIDs"`
	}
	params := controller.Validated[editStaffParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	request := corporationdto.EditStaffRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		StaffID:       params.StaffID,
		Status:        params.Status,
		RoleIDs:       params.RoleIDs,
	}
	if err := corporationController.corporationService.EditStaff(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.editStaff")
	controller.Response(ctx, 200, message, nil)
}
