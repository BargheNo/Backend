package maintenance

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/gin-gonic/gin"
)

type CustomerMaintenanceController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
}

func NewCustomerMaintenanceController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
) *CustomerMaintenanceController {
	return &CustomerMaintenanceController{
		constants:  constants,
		pagination: pagination,
	}
}

func (maintenanceController *CustomerMaintenanceController) CreateMaintenanceRequest(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		PanelID                 uint   `json:"panelID" validate:"required"`
		Subject                 string `json:"subject" validate:"required"`
		Description             string `json:"description" validate:"required"`
		SameInstallationCompany bool   `json:"sameInstallationCompany"`
		UrgencyLevel            string `json:"urgencyLevel" validate:"required"`
	}
}
