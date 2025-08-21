package rbac

import (
	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationRBACController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	rbacService usecase.RBACService
}

func NewCorporationRBACController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	rbacService usecase.RBACService,
) *CorporationRBACController {
	return &CorporationRBACController{
		constants:   constants,
		pagination:  pagination,
		rbacService: rbacService,
	}
}

func (corporationController *CorporationRBACController) GetCorporationRoles(ctx *gin.Context) {
	type getRolesParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	params := controller.Validated[getRolesParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)

	request := rbacdto.GetRolesListRequest{
		IsStaff: true,
		Query:   params.Query,
		Offset:  offset,
		Limit:   limit,
	}

	roles, count, err := corporationController.rbacService.GetRoles(request)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(roles, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
