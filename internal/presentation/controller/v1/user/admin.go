package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	constants   *bootstrap.Constants
	userService service.UserService
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	userService service.UserService,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		userService: userService,
	}
}

func (userController *AdminUserController) GetPermissionsList(ctx *gin.Context) {
	permissions := userController.userService.GetAllPermissions()
	controller.Response(ctx, 200, "", permissions)
}

func (userController *AdminUserController) GetRolesList(ctx *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) CreateRole(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) GetRoleDetails(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) GetRoleOwners(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) UpdateRole(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) DeleteRole(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) GetUserRoles(c *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) UpdateUserRoles(c *gin.Context) {
	// some codes here ...
}
