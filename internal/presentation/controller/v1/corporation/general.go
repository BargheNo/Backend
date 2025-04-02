package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
)

type GeneralCorporationController struct {
	constants          *bootstrap.Constants
	corporationService service.CorporationService
}

func NewGeneralCorporationController(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
) *GeneralCorporationController {
	return &GeneralCorporationController{
		constants:          constants,
		corporationService: corporationService,
	}
}
