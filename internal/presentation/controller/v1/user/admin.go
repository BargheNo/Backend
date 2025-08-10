package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	userService usecase.UserService
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	userService usecase.UserService,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		pagination:  pagination,
		userService: userService,
	}
}

func (userController *AdminUserController) GetPermissionsList(ctx *gin.Context) {
	type getPermissionsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := controller.Validated[getPermissionsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetPermissionsListRequest{
		Offset: offset,
		Limit:  limit,
	}

	permissions, count, err := userController.userService.GetAllPermissions(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(permissions, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetPermissionRoles(ctx *gin.Context) {
	type getPermissionRolesParams struct {
		PermissionID uint `uri:"permissionID" validate:"required"`
		Page         int  `form:"page"`
		PageSize     int  `form:"pageSize"`
		SortBy       uint `form:"sortBy"`
		Asc          bool `form:"asc"`
	}
	params := controller.Validated[getPermissionRolesParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetPermissionRolesRequest{
		PermissionID: params.PermissionID,
		Offset:       offset,
		Limit:        limit,
		SortBy:       params.SortBy,
		Asc:          params.Asc,
	}

	roles, count, err := userController.userService.GetPermissionRoles(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roles, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetRolesList(ctx *gin.Context) {
	type getRolesParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := controller.Validated[getRolesParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetRolesListRequest{
		Offset: offset,
		Limit:  limit,
	}

	roles, count, err := userController.userService.GetAllRoles(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roles, count, offset, limit)
	controller.Response(ctx, 200, "", data)
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

	role, err := userController.userService.GetRoleDetails(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", role)
}

func (userController *AdminUserController) GetRoleOwners(ctx *gin.Context) {
	type getRoleParams struct {
		RoleID   uint `uri:"roleID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	params := controller.Validated[getRoleParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetRoleOwnersRequest{
		RoleID: params.RoleID,
		Offset: offset,
		Limit:  limit,
	}

	roleOwners, count, err := userController.userService.GetRoleOwners(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roleOwners, count, offset, limit)
	controller.Response(ctx, 200, "", data)
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
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[usersParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetUsersListRequest{
		Status: params.Status,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	users, count, err := userController.userService.GetUsersByStatus(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(users, count, offset, limit)

	controller.Response(ctx, 200, "", data)
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
