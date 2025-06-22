package payment

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralPaymentController struct {
	constants      *bootstrap.Constants
	paymentService port.PaymentService
}

func NewGeneralPaymentController(
	constants *bootstrap.Constants,
	paymentService port.PaymentService,
) *GeneralPaymentController {
	return &GeneralPaymentController{
		constants:      constants,
		paymentService: paymentService,
	}
}

func (corporationController *GeneralPaymentController) GetPaymentMethods(ctx *gin.Context) {
	methods := corporationController.paymentService.GetPaymentMethods()
	controller.Response(ctx, 200, "", methods)
}
