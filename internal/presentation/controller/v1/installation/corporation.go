package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService service.InstallationService
}

func NewCorporationInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService service.InstallationService,
) *CorporationInstallationController {
	return &CorporationInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *CorporationInstallationController) GetInstallationRequests(ctx *gin.Context) {
	type getInstallationRequestParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getInstallationRequestParams](ctx)
	operatorID, _ := ctx.Get(installationController.constants.Context.ID)

	pagination := controller.GetPagination(ctx, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := installationdto.CorporationPanelListRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		Offset:        offset,
		Limit:         limit,
	}
	installationRequest, err := installationController.installationService.GetAnonymousInstallationRequests(listInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", installationRequest)
}

func (installationController *CorporationInstallationController) GetInstallationRequest(ctx *gin.Context) {
	type getInstallationRequestParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		RequestID     uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[getInstallationRequestParams](ctx)
	operatorID, _ := ctx.Get(installationController.constants.Context.ID)

	requestInfo := installationdto.CorporationPanelRequest{
		CorporationID:  params.CorporationID,
		OperatorID:     operatorID.(uint),
		InstallationID: params.RequestID,
	}
	installationRequest, err := installationController.installationService.GetAnonymousInstallationRequest(requestInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", installationRequest)
}

func (installationController *CorporationInstallationController) CompleteInstallation(ctx *gin.Context) {
	type completeBidParams struct {
		CorporationID   uint `uri:"corporationID" validate:"required"`
		PanelID         uint `uri:"panelID" validate:"required"`
		Tilt            uint `json:"tilt" validate:"required"`
		Azimuth         uint `json:"area" validate:"required"`
		NumberOfModules uint `json:"numberOfModules" validate:"required"`
	}
	params := controller.Validated[completeBidParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	panelInfo := installationdto.CompleteBidInstallationRequest{
		CorporationID:   params.CorporationID,
		OperatorID:      userID.(uint),
		PanelID:         params.PanelID,
		Tilt:            params.Tilt,
		Azimuth:         params.Azimuth,
		NumberOfModules: params.NumberOfModules,
	}
	err := installationController.installationService.CompleteInstallationRequest(panelInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.completeInstallation")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CorporationInstallationController) AddPanel(ctx *gin.Context) {
	type addPanelParams struct {
		CorporationID        uint   `uri:"corporationID" validate:"required"`
		Name                 string `json:"name" validate:"required"`
		CustomerPhone        string `json:"customerPhone" validate:"required"`
		Power                uint   `json:"power" validate:"required"`
		Area                 uint   `json:"area" validate:"required"`
		BuildingType         uint   `json:"buildingType" validate:"required"`
		Tilt                 uint   `json:"tilt" validate:"required"`
		Azimuth              uint   `json:"azimuth" validate:"required"`
		TotalNumberOfModules uint   `json:"totalNumberOfModules" validate:"required"`
		ProvinceID           uint   `json:"provinceID" validate:"required"`
		CityID               uint   `json:"cityID" validate:"required"`
		StreetAddress        string `json:"streetAddress" validate:"required"`
		PostalCode           string `json:"postalCode" validate:"required"`
		HouseNumber          string `json:"houseNumber" validate:"required"`
		Unit                 uint   `json:"unit" validate:"required"`
		GuaranteeID          *uint  `json:"guaranteeID"`
	}
	params := controller.Validated[addPanelParams](ctx)
	operatorID, _ := ctx.Get(installationController.constants.Context.ID)

	panelInfo := installationdto.AddPanelRequest{
		Name:                 params.Name,
		CorporationID:        params.CorporationID,
		OperatorID:           operatorID.(uint),
		CustomerPhone:        params.CustomerPhone,
		Status:               enum.PanelStatusActive,
		Power:                params.Power,
		Area:                 params.Area,
		BuildingType:         params.BuildingType,
		Tilt:                 params.Tilt,
		Azimuth:              params.Azimuth,
		GuaranteeID:          params.GuaranteeID,
		TotalNumberOfModules: params.TotalNumberOfModules,
		Address: addressdto.CreateAddressRequest{
			ProvinceID:    params.ProvinceID,
			CityID:        params.CityID,
			StreetAddress: params.StreetAddress,
			PostalCode:    params.PostalCode,
			HouseNumber:   params.HouseNumber,
			Unit:          params.Unit,
		},
	}
	err := installationController.installationService.AddPanel(panelInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addPanel")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *CorporationInstallationController) GetCorporationPanels(ctx *gin.Context) {
	type getInstallationRequestParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		Status        uint `form:"status" validate:"required"`
	}
	params := controller.Validated[getInstallationRequestParams](ctx)
	operatorID, _ := ctx.Get(installationController.constants.Context.ID)

	pagination := controller.GetPagination(ctx, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := installationdto.CorporationPanelListRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		Status:        params.Status,
		Offset:        offset,
		Limit:         limit,
	}
	panels, err := installationController.installationService.GetCorporationPanels(listInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", panels)
}

func (installationController *CorporationInstallationController) GetCorporationPanel(ctx *gin.Context) {
	type getPanelParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[getPanelParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	panelInfo := installationdto.CorporationPanelRequest{
		CorporationID:  params.CorporationID,
		OperatorID:     userID.(uint),
		InstallationID: params.PanelID,
	}
	panel, err := installationController.installationService.GetCorporationPanel(panelInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", panel)
}

func (installationController *CorporationInstallationController) ViolatePanelGuarantee(ctx *gin.Context) {
	type panelGuaranteeViolationParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		PanelID       uint   `uri:"panelID" validate:"required"`
		Reason        string `json:"reason" validate:"required"`
		Details       string `json:"details" validate:"required"`
	}
	params := controller.Validated[panelGuaranteeViolationParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	violationInfo := installationdto.CreateViolatePanelGuaranteeRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		PanelID:       params.PanelID,
		GuaranteeViolation: guaranteedto.CreateGuaranteeViolationRequest{
			PanelID:       params.PanelID,
			CorporationID: params.CorporationID,
			OperatorID:    userID.(uint),
			Reason:        params.Reason,
			Details:       params.Details,
		},
	}
	_, err := installationController.installationService.ViolatePanelGuaranteeStatus(violationInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addGuaranteeViolation")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CorporationInstallationController) ClearPanelGuaranteeViolation(ctx *gin.Context) {
	type panelGuaranteeViolationParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[panelGuaranteeViolationParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	violationInfo := installationdto.GetCorporationGuaranteeViolationRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		PanelID:       params.PanelID,
	}
	err := installationController.installationService.ClearPanelGuaranteeViolation(violationInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.clearGuaranteeViolation")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CorporationInstallationController) GetPanelGuaranteeViolation(ctx *gin.Context) {
	type updateGuaranteeParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[updateGuaranteeParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	request := installationdto.GetCorporationGuaranteeViolationRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		PanelID:       params.PanelID,
	}
	violation, err := installationController.installationService.GetCorporationPanelGuaranteeViolation(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", violation)
}

func (installationController *CorporationInstallationController) UpdatePanelGuaranteeViolation(ctx *gin.Context) {
	type updateGuaranteeParams struct {
		CorporationID uint    `uri:"corporationID" validate:"required"`
		PanelID       uint    `uri:"panelID" validate:"required"`
		Reason        *string `json:"reason"`
		Details       *string `json:"details"`
	}
	params := controller.Validated[updateGuaranteeParams](ctx)
	userID, _ := ctx.Get(installationController.constants.Context.ID)

	request := installationdto.UpdateGuaranteeViolationRequest{
		CorporationID: params.CorporationID,
		PanelID:       params.PanelID,
		OperatorID:    userID.(uint),
		Reason:        params.Reason,
		Details:       params.Details,
	}
	err := installationController.installationService.UpdatePanelGuaranteeViolation(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateGuaranteeViolation")
	controller.Response(ctx, 200, message, nil)
}
