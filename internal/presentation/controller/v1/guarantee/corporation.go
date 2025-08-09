package guarantee

import (
	"github.com/BargheNo/Backend/bootstrap"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationGuaranteeController struct {
	constants        *bootstrap.Constants
	guaranteeService usecase.GuaranteeService
}

func NewCorporationGuaranteeController(
	constants *bootstrap.Constants,
	guaranteeService usecase.GuaranteeService,
) *CorporationGuaranteeController {
	return &CorporationGuaranteeController{
		constants:        constants,
		guaranteeService: guaranteeService,
	}
}

func (guaranteeController *CorporationGuaranteeController) CreateGuarantee(ctx *gin.Context) {
	type guaranteeTerms struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Limitations string `json:"limitations"`
	}
	type createGuaranteeParams struct {
		CorporationID  uint             `uri:"corporationID" validate:"required"`
		Name           string           `json:"name" validate:"required"`
		GuaranteeType  uint             `json:"type" validate:"required"`
		Duration       uint             `json:"duration" validate:"required"`
		Description    string           `json:"description"`
		GuaranteeTerms []guaranteeTerms `json:"terms" validate:"required"`
	}
	params := controller.Validated[createGuaranteeParams](ctx)
	userID, _ := ctx.Get(guaranteeController.constants.Context.ID)

	guaranteeTermsRequest := make([]guaranteedto.GuaranteeTermsRequest, len(params.GuaranteeTerms))
	for i, plan := range params.GuaranteeTerms {
		guaranteeTermsRequest[i] = guaranteedto.GuaranteeTermsRequest{
			Title:       plan.Title,
			Description: plan.Description,
			Limitations: plan.Limitations,
		}
	}

	request := guaranteedto.CreateGuaranteeRequest{
		CorporationID:         params.CorporationID,
		OperatorID:            userID.(uint),
		Name:                  params.Name,
		Status:                enum.GuaranteeStatusActive,
		GuaranteeType:         params.GuaranteeType,
		Duration:              params.Duration,
		Description:           params.Description,
		GuaranteeTermsRequest: guaranteeTermsRequest,
	}
	guaranteeID, err := guaranteeController.guaranteeService.AddGuarantee(request)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, guaranteeController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createGuarantee")
	controller.Response(ctx, 200, message, guaranteeID)
}

func (guaranteeController *CorporationGuaranteeController) GetGuaranteeTypes(ctx *gin.Context) {
	guaranteeTypes := guaranteeController.guaranteeService.GetGuaranteeTypes()
	controller.Response(ctx, 200, "", guaranteeTypes)
}

func (guaranteeController *CorporationGuaranteeController) GetGuaranteeStatuses(ctx *gin.Context) {
	guaranteeTypes := guaranteeController.guaranteeService.GetGuaranteeStatuses()
	controller.Response(ctx, 200, "", guaranteeTypes)
}

func (guaranteeController *CorporationGuaranteeController) GetGuarantees(ctx *gin.Context) {
	type getGuaranteesParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		Status        uint `form:"status"`
	}
	params := controller.Validated[getGuaranteesParams](ctx)
	userID, _ := ctx.Get(guaranteeController.constants.Context.ID)

	request := guaranteedto.GetGuaranteesRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		Status:        params.Status,
	}
	guarantees, err := guaranteeController.guaranteeService.GetCorporationGuarantees(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", guarantees)
}

func (guaranteeController *CorporationGuaranteeController) GetGuarantee(ctx *gin.Context) {
	type getGuaranteeParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		GuaranteeID   uint `uri:"guaranteeID" validate:"required"`
	}
	params := controller.Validated[getGuaranteeParams](ctx)
	userID, _ := ctx.Get(guaranteeController.constants.Context.ID)

	request := guaranteedto.GetGuaranteeRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		GuaranteeID:   params.GuaranteeID,
	}
	guarantees, err := guaranteeController.guaranteeService.GetCorporationGuarantee(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", guarantees)
}

func (guaranteeController *CorporationGuaranteeController) UpdateGuarantee(ctx *gin.Context) {
	type updateGuaranteeParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		GuaranteeID   uint `uri:"guaranteeID" validate:"required"`
		Status        uint `json:"status" validate:"required"`
	}
	params := controller.Validated[updateGuaranteeParams](ctx)
	userID, _ := ctx.Get(guaranteeController.constants.Context.ID)

	request := guaranteedto.ChangeStatusRequest{
		CorporationID: params.CorporationID,
		OperatorID:    userID.(uint),
		GuaranteeID:   params.GuaranteeID,
		Status:        params.Status,
	}
	if err := guaranteeController.guaranteeService.UpdateGuaranteeStatus(request); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, guaranteeController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateGuarantee")
	controller.Response(ctx, 200, message, nil)
}
