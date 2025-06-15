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
	pagination  *bootstrap.Pagination
	userService service.UserService
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	userService service.UserService,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		pagination:  pagination,
		userService: userService,
	}
}

func (userController *AdminUserController) GetPermissionsList(ctx *gin.Context) {
	permissions, err := userController.userService.GetAllPermissions()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", permissions)
}

func (userController *AdminUserController) GetRolesList(ctx *gin.Context) {
	roles, err := userController.userService.GetAllRoles()
	if err != nil {
		panic(err)
	}
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
	if err := userController.userService.CreateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetRoleDetails(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[getRoleParams](ctx)

	role, err := userController.userService.GetRoomDetails(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", role)
}

func (userController *AdminUserController) GetRoleOwners(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[getRoleParams](ctx)
	roleOwners, err := userController.userService.GetRoleOwners(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roleOwners)
}

func (userController *AdminUserController) UpdateRole(ctx *gin.Context) {
	type updateRoleParams struct {
		RoleID        uint    `uri:"roleID" validate:"required"`
		Name          *string `json:"name"`
		PermissionIDs []uint  `json:"permissionIDs"`
	}
	params := controller.Validated[updateRoleParams](ctx)

	newRoleRequest := userdto.UpdateRoleRequest{
		RoleID:        params.RoleID,
		Name:          params.Name,
		PermissionIDs: params.PermissionIDs,
	}
	if err := userController.userService.UpdateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) DeleteRole(ctx *gin.Context) {
	type deleteRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[deleteRoleParams](ctx)

	if err := userController.userService.DeleteRole(params.RoleID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteRole")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetUserRoles(ctx *gin.Context) {
	type getRolesParams struct {
		UserID uint `uri:"userID" validate:"required"`
	}
	params := controller.Validated[getRolesParams](ctx)
	roles, err := userController.userService.GetUserRoles(params.UserID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roles)
}

func (userController *AdminUserController) UpdateUserRoles(ctx *gin.Context) {
	type updateUserRolesParams struct {
		UserID  uint   `uri:"userID" validate:"required"`
		RoleIDs []uint `json:"roleIDs"`
	}
	params := controller.Validated[updateUserRolesParams](ctx)

	userRolesRequest := userdto.UpdateUserRolesRequest{
		UserID:  params.UserID,
		RoleIDs: params.RoleIDs,
	}
	if err := userController.userService.UpdateUserRoles(userRolesRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateUserRoles")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) GetUsers(ctx *gin.Context) {
	type usersParams struct {
		Statuses []uint `form:"statuses"`
	}
	params := controller.Validated[usersParams](ctx)
	pagination := controller.GetPagination(ctx, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	request := userdto.GetUsersListRequest{
		Statuses: params.Statuses,
		Offset:   offset,
		Limit:    limit,
	}
	users, err := userController.userService.GetUsersByStatus(request)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", users)
}

func (userController *AdminUserController) BanUser(ctx *gin.Context) {
	type banParams struct {
		UserID uint `uri:"userID"`
	}
	params := controller.Validated[banParams](ctx)

	if err := userController.userService.BanUser(params.UserID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.banUser")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) UnbanUser(ctx *gin.Context) {
	type unbanParams struct {
		UserID uint `uri:"userID"`
	}
	params := controller.Validated[unbanParams](ctx)

	if err := userController.userService.UnbanUser(params.UserID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unbanUser")
	controller.Response(ctx, 200, message, nil)
}
