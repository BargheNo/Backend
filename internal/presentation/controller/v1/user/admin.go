package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
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
	roles := userController.userService.GetAllRoles()
	controller.Response(ctx, 200, "", roles)
}

func (userController *AdminUserController) CreateRole(ctx *gin.Context) {
	type newRoleParams struct {
		Name          string `json:"name" validate:"required"`
		PermissionIDs []uint `json:"permissionIDs"`
	}
	params := controller.Validated[newRoleParams](ctx)

	newRoleRequest := userdto.NewRoleRequest{
		Name:          params.Name,
		PermissionIDs: params.PermissionIDs,
	}
	userController.userService.CreateRole(newRoleRequest)

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetRoleDetails(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[getRoleParams](ctx)

	role := userController.userService.GetRoomDetails(params.RoleID)
	controller.Response(ctx, 200, "", role)
}

func (userController *AdminUserController) GetRoleOwners(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[getRoleParams](ctx)
	roleOwners := userController.userService.GetRoleOwners(param.RoleID)
	controller.Response(ctx, 200, "", roleOwners)
}

func (userController *AdminUserController) UpdateRole(ctx *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) DeleteRole(ctx *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) GetUserRoles(ctx *gin.Context) {
	// some codes here ...
}

func (userController *AdminUserController) UpdateUserRoles(ctx *gin.Context) {
	// some codes here ...
}
