package rbac

import (
	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminRBACController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	rbacService usecase.RBACService
}

func NewAdminRBACController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	rbacService usecase.RBACService,
) *AdminRBACController {
	return &AdminRBACController{
		constants:   constants,
		pagination:  pagination,
		rbacService: rbacService,
	}
}

func (rbacController *AdminRBACController) GetPermissionsList(ctx *gin.Context) {
	type getPermissionsParams struct {
		IsStaff  bool `form:"isStaff"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	params := controller.Validated[getPermissionsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, rbacController.pagination.DefaultPage, rbacController.pagination.DefaultPageSize)

	request := rbacdto.GetPermissionsListRequest{
		IsStaff: params.IsStaff,
		Offset:  offset,
		Limit:   limit,
	}

	permissions, count, err := rbacController.rbacService.GetPermissions(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(permissions, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (rbacController *AdminRBACController) GetPermissionRoles(ctx *gin.Context) {
	type getPermissionRolesParams struct {
		PermissionID uint `uri:"permissionID" validate:"required"`
		Page         int  `form:"page"`
		PageSize     int  `form:"pageSize"`
		SortBy       uint `form:"sortBy"`
		Asc          bool `form:"asc"`
	}
	params := controller.Validated[getPermissionRolesParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, rbacController.pagination.DefaultPage, rbacController.pagination.DefaultPageSize)

	request := rbacdto.GetPermissionRolesRequest{
		PermissionID: params.PermissionID,
		Offset:       offset,
		Limit:        limit,
		SortBy:       params.SortBy,
		Asc:          params.Asc,
	}

	roles, count, err := rbacController.rbacService.GetPermissionRoles(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roles, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (rbacController *AdminRBACController) GetRolesList(ctx *gin.Context) {
	type getRolesParams struct {
		IsStaff  bool   `form:"isStaff"`
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	params := controller.Validated[getRolesParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, rbacController.pagination.DefaultPage, rbacController.pagination.DefaultPageSize)

	request := rbacdto.GetRolesListRequest{
		IsStaff: params.IsStaff,
		Query:   params.Query,
		Offset:  offset,
		Limit:   limit,
	}

	roles, count, err := rbacController.rbacService.GetRoles(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roles, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (rbacController *AdminRBACController) CreateRole(ctx *gin.Context) {
	type createRoleParams struct {
		Name          string `json:"name" validate:"required"`
		IsStaff       bool   `json:"isStaff"`
		PermissionIDs []uint `json:"permissionIDs"`
	}
	params := controller.Validated[createRoleParams](ctx)

	newRoleRequest := rbacdto.NewRoleRequest{
		Name:          params.Name,
		IsStaff:       params.IsStaff,
		PermissionIDs: params.PermissionIDs,
	}

	if err := rbacController.rbacService.CreateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, rbacController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createRole")
	controller.Response(ctx, 200, message, nil)
}

func (rbacController *AdminRBACController) GetRoleDetails(ctx *gin.Context) {
	type getRoleDetailsParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[getRoleDetailsParams](ctx)

	role, err := rbacController.rbacService.GetRoleDetails(params.RoleID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", role)
}

func (rbacController *AdminRBACController) UpdateRole(ctx *gin.Context) {
	type updateRoleParams struct {
		RoleID        uint    `uri:"roleID" validate:"required"`
		Name          *string `json:"name"`
		PermissionIDs []uint  `json:"permissionIDs"`
	}
	params := controller.Validated[updateRoleParams](ctx)

	newRoleRequest := rbacdto.UpdateRoleRequest{
		RoleID:        params.RoleID,
		Name:          params.Name,
		PermissionIDs: params.PermissionIDs,
	}
	if err := rbacController.rbacService.UpdateRole(newRoleRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, rbacController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateRole")
	controller.Response(ctx, 200, message, nil)
}

func (rbacController *AdminRBACController) DeleteRole(ctx *gin.Context) {
	type deleteRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	params := controller.Validated[deleteRoleParams](ctx)

	if err := rbacController.rbacService.DeleteRole(params.RoleID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, rbacController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteRole")
	controller.Response(ctx, 200, message, nil)
}
