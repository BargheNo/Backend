package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService service.CorporationService
}

func NewAdminCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService service.CorporationService,
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
		Status uint `form:"status" validate:"required"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	pagination := controller.GetPagination(ctx, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := corporationdto.GetCorporationsByAdminRequest{
		Status: params.Status,
		Limit:  limit,
		Offset: offset,
	}
	corporations := corporationController.corporationService.GetCorporationsByAdmin(listInfo)

	controller.Response(ctx, 200, "", corporations)
}

func (corporationController *AdminCorporationController) GetCorporation(ctx *gin.Context) {
	type getCorporationsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationsParams](ctx)

	corporation := corporationController.corporationService.GetCorporationByAdmin(params.CorporationID)

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

	reviews := corporationController.corporationService.GetCorporationReviewsByAdmin(params.CorporationID)

	controller.Response(ctx, 200, "", reviews)
}

func (corporationController *AdminCorporationController) ApproveCorporationRequest(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) RejectCorporationRequest(ctx *gin.Context) {
	// some codes here ...
}
