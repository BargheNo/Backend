package httpv1

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	ticket := routerGroup.Group("/ticket")
	{
		ticket.GET("", app.Controllers.Admin.TicketController.GetTickets)
		ticket.GET("/:ticketID/comments", app.Controllers.Admin.TicketController.GetComments)
		ticket.POST("/:ticketID/comments", app.Controllers.Admin.TicketController.CreateComment)
		ticket.POST("/:ticketID/resolve", app.Controllers.Admin.TicketController.ResolveTicket)
	}

	accessManagement := routerGroup.Group("")
	// accessManagement.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManageUsers, enums.ManageRoles}))
	{
		accessManagement.GET("/permissions", app.Controllers.Admin.UserController.GetPermissionsList)

		roles := accessManagement.Group("/roles")
		{
			roles.GET("", app.Controllers.Admin.UserController.GetRolesList)
			roles.POST("", app.Controllers.Admin.UserController.CreateRole)

			rolesSubGroup := roles.Group("/:roleID")
			{
				rolesSubGroup.GET("", app.Controllers.Admin.UserController.GetRoleDetails)
				rolesSubGroup.GET("/owners", app.Controllers.Admin.UserController.GetRoleOwners)
				rolesSubGroup.PUT("", app.Controllers.Admin.UserController.UpdateRole)
				rolesSubGroup.DELETE("", app.Controllers.Admin.UserController.DeleteRole)
			}
		}

		userRoles := accessManagement.Group("/users/:userID/roles")
		{
			userRoles.GET("", app.Controllers.Admin.UserController.GetUserRoles)
			userRoles.PUT("", app.Controllers.Admin.UserController.UpdateUserRoles)
		}
	}

	userManagement := routerGroup.Group("/users")
	userManagement.Use(app.Middlewares.Authentication.RequiredWithPermission([]enum.PermissionType{enum.UserManageRolesPermissions}))
	{
		userManagement.GET("", app.Controllers.Admin.UserController.GetUsers)
		userManagement.PUT("/:userID/ban", app.Controllers.Admin.UserController.BanUser)
		userManagement.PUT("/:userID/unban", app.Controllers.Admin.UserController.UnbanUser)
	}

	report := routerGroup.Group("/report")
	{
		report.GET("/maintenance", app.Controllers.Admin.ReportController.GetMaintenanceReports)
		report.GET("/panel", app.Controllers.Admin.ReportController.GetPanelReports)
		report.POST("/resolve/:reportID", app.Controllers.Admin.ReportController.ResolveReport)
	}

	news := routerGroup.Group("/news")
	{
		news.POST("/draft")
		news.GET("")

		newsSubgroup := news.Group("/:newsID")
		{
			newsSubgroup.POST("/finalize")
			newsSubgroup.GET("")
			newsSubgroup.PUT("")
			newsSubgroup.DELETE("")
		}

	}
}
