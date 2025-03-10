package sample

import (
	"github.com/BargheNo/Backend/bootstrap"
	service_interfaces "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type SampleController struct {
	constants     *bootstrap.Constants
	sampleService service_interfaces.SampleService
}

func NewSampleController(
	constants *bootstrap.Constants,
	sampleService service_interfaces.SampleService,
) *SampleController {
	return &SampleController{
		constants:     constants,
		sampleService: sampleService,
	}
}

func (sampleController *SampleController) SampleCreate(c *gin.Context) {
	sampleController.sampleService.SampleCreate()

	trans := controller.GetTranslator(c, sampleController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.userRegistration")
	controller.Response(c, 200, message, nil)
}
