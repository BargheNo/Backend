package user

import (
	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
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

func (userController *AdminUserController) GetUsers(ctx *gin.Context) {
	type usersParams struct {
		Query    string `form:"query"`
		Status   uint   `form:"status"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[usersParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := userdto.GetUsersListRequest{
		Query:  params.Query,
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

func (userController *AdminUserController) GetRoleOwners(ctx *gin.Context) {
	type getRoleOwnersParams struct {
		RoleID   uint   `uri:"roleID" validate:"required"`
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[getRoleOwnersParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, userController.pagination.DefaultPage, userController.pagination.DefaultPageSize)

	request := rbacdto.GetRoleOwnersRequest{
		RoleID: params.RoleID,
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	users, count, err := userController.userService.GetRoleOwners(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(users, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
