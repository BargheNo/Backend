package httpv1

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	const status string = "/status"

	ticket := routerGroup.Group("/ticket")
	{
		ticket.GET("", app.Controllers.Admin.TicketController.GetTickets)
		ticket.GET("/:ticketID/comments", app.Controllers.Admin.TicketController.GetComments)
		ticket.POST("/:ticketID/comments", app.Controllers.Admin.TicketController.CreateComment)
		ticket.POST("/:ticketID/resolve", app.Controllers.Admin.TicketController.ResolveTicket)
	}

	installations := routerGroup.Group("/installation")
	{
		requests := installations.Group("/request")
		{
			requests.GET("", app.Controllers.Admin.InstallationController.GetInstallationRequests)
			requestSubGroup := requests.Group("/:requestID")
			{
				requestSubGroup.GET("", app.Controllers.Admin.InstallationController.GetInstallationRequest)
				requestSubGroup.DELETE("", app.Controllers.Admin.InstallationController.DeleteInstallationRequest)
				requestSubGroup.PUT("", app.Controllers.Admin.InstallationController.UpdateInstallationRequest)
				requestSubGroup.GET("/bid", app.Controllers.Admin.BidController.GetBids)
			}
		}

		panels := installations.Group("/panel")
		{
			panels.GET("", app.Controllers.Admin.InstallationController.GetPanels)
			panels.GET(status, app.Controllers.Admin.InstallationController.GetAllPanelStatuses)
			panelsSubGroup := panels.Group("/:panelID")
			{
				panelsSubGroup.GET("", app.Controllers.Admin.InstallationController.GetPanel)
				panelsSubGroup.PUT("", app.Controllers.Admin.InstallationController.UpdatePanel)
				panelsSubGroup.DELETE("", app.Controllers.Admin.InstallationController.DeletePanel)
			}
		}
	}

	bids := routerGroup.Group("/bids")
	{
		bids.GET("")
		bidsSubGroup := bids.Group("/:bidID")
		{
			bidsSubGroup.GET("")
			bidsSubGroup.PUT("")
			bidsSubGroup.DELETE("")
		}
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

	corporationManagement := routerGroup.Group("/corporation")
	// put some role here
	{
		corporationManagement.GET("", app.Controllers.Admin.CorporationController.GetCorporations)
		corporationManagement.GET(status, app.Controllers.Admin.CorporationController.GetCorporationStatus)
		corporationSubGRoup := corporationManagement.Group("/:corporationID")
		{
			corporationSubGRoup.GET("", app.Controllers.Admin.CorporationController.GetCorporation)
			corporationSubGRoup.POST("/approve", app.Controllers.Admin.CorporationController.ApproveCorporationRequest)
			corporationSubGRoup.POST("/reject", app.Controllers.Admin.CorporationController.RejectCorporationRequest)
			reviewSubGroup := corporationSubGRoup.Group("/review")
			{
				reviewSubGroup.GET("/action", app.Controllers.Admin.CorporationController.GetReviewActions)
				reviewSubGroup.GET("", app.Controllers.Admin.CorporationController.GetCorporationReviews)
			}
		}
	}

	report := routerGroup.Group("/report")
	{
		report.GET("/maintenance", app.Controllers.Admin.ReportController.GetMaintenanceReports)
		report.GET("/panel", app.Controllers.Admin.ReportController.GetPanelReports)
		report.POST("/resolve/:reportID", app.Controllers.Admin.ReportController.ResolveReport)
	}

	news := routerGroup.Group("/news")
	{
		news.POST("/draft", app.Controllers.Admin.NewsController.CreateDraftNews)
		news.GET("", app.Controllers.Admin.NewsController.GetNewsList)
		news.GET(status, app.Controllers.Admin.NewsController.GetAllNewsStatuses)
		news.DELETE("", app.Controllers.Admin.NewsController.DeleteNews)
		newsSubgroup := news.Group("/:newsID")
		{
			newsSubgroup.GET("", app.Controllers.Admin.NewsController.GetNews)
			newsSubgroup.PUT("", app.Controllers.Admin.NewsController.EditNews)
			newsSubgroup.PUT("/publish", app.Controllers.Admin.NewsController.PublishNews)
			newsSubgroup.PUT("unpublish", app.Controllers.Admin.NewsController.UnpublishNews)
			newsSubgroup.POST("/media", app.Controllers.Admin.NewsController.AddNewsMedia)
			newsSubgroup.DELETE("/media/:mediaID", app.Controllers.Admin.NewsController.DeleteNewsMedia)
			newsSubgroup.GET("/media/:mediaID", app.Controllers.Admin.NewsController.GetNewsMedia)
		}

	}
}
