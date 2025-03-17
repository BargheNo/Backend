package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
)

type CustomerUserController struct {
	constants   *bootstrap.Constants
	userService service.UserService
}

func NewCustomerUserController(
	constants *bootstrap.Constants,
	userService service.UserService,
) *CustomerUserController {
	return &CustomerUserController{
		constants:   constants,
		userService: userService,
	}
}
