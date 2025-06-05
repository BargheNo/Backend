package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
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

func (corporationController *AdminCorporationController) GetCorporations(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) GetCorporationStatus(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) GetCorporation(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) GetCorporationReview(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) ApproveCorporationRequest(ctx *gin.Context) {
	// some codes here ...
}

func (corporationController *AdminCorporationController) RejectCorporationRequest(ctx *gin.Context) {
	// some codes here ...
}
