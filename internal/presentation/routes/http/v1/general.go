package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	const (
		status   = "/status"
		sortable = "/sortable"
	)

	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register/basic", app.Controllers.General.UserController.BasicRegister)
		auth.POST("/verify/phone", app.Controllers.General.UserController.VerifyPhone)
		auth.POST("/login", app.Controllers.General.UserController.Login)
		auth.POST("/forgot-password", app.Controllers.General.UserController.ForgotPassword)
		auth.POST("/confirm-otp", app.Controllers.General.UserController.ConfirmOTP)
		auth.POST("/refresh", app.Controllers.General.UserController.RefreshToken)
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.GET("/province", app.Controllers.General.AddressController.GetProvince)
		addresses.GET("/province/:provinceID/city", app.Controllers.General.AddressController.GetProvinceCities)
	}

	contacts := routerGroup.Group("/contact")
	{
		contacts.GET("/types", app.Controllers.General.CorporationController.GetContactTypes)
	}

	corporations := routerGroup.Group("/corporation")
	{
		corporations.GET("", app.Controllers.General.CorporationController.GetCorporations)
		corporations.GET(sortable, app.Controllers.General.CorporationController.GetSortableFields)
		corporations.GET("staff/sortable", app.Controllers.General.CorporationController.GetStaffSortableFields)
	}

	notifications := routerGroup.Group("/notifications")
	{
		notifications.GET("/type", app.Controllers.General.NotificationController.GetContactTypes)
		notifications.GET(sortable, app.Controllers.General.NotificationController.GetSortableFields)
	}

	installations := routerGroup.Group("/installation")
	{
		requests := installations.Group("/request")
		{
			requests.GET(status, app.Controllers.General.InstallationController.GetRequestStatuses)
			requests.GET("/building", app.Controllers.General.InstallationController.GetBuildingTypes)
			requests.GET(sortable, app.Controllers.General.InstallationController.GetRequestSortableFields)
		}
		panels := installations.Group("/panel")
		{
			panels.GET(status, app.Controllers.General.InstallationController.GetPanelStatuses)
			panels.GET(sortable, app.Controllers.General.InstallationController.GetPanelSortableFields)
		}
	}

	guarantees := routerGroup.Group("/guarantee")
	{
		guarantees.GET(status, app.Controllers.Corporation.GuaranteeController.GetGuaranteeStatuses)
	}

	maintenances := routerGroup.Group("/maintenance")
	{
		maintenances.GET(status, app.Controllers.Customer.MaintenanceController.GetMaintenanceStatuses)
		maintenances.GET(sortable, app.Controllers.General.MaintenanceController.GetSortableFields)
	}

	payments := routerGroup.Group("/payment")
	{
		payments.GET("method", app.Controllers.General.PaymentController.GetPaymentMethods)
	}

	news := routerGroup.Group("/news")
	{
		news.GET("", app.Controllers.General.NewsController.GetNewsList)
		news.GET("/:newsID", app.Controllers.General.NewsController.GetNews)
		news.GET("/:newsID/media/:mediaID", app.Controllers.General.NewsController.GetNewsMedia)
		news.GET(sortable, app.Controllers.General.NewsController.GetSortableFields)
	}

	blogs := routerGroup.Group("/blog")
	{
		blogs.GET("", app.Controllers.General.BlogController.GetPosts)
		blogs.GET("/corporation/:corporationID", app.Controllers.General.BlogController.GetCorporationPosts)
		blogs.GET("/:postID", app.Controllers.General.BlogController.GetPost)
		blogs.GET("/:postID/media/:mediaID", app.Controllers.General.BlogController.GetPostMedia)
		blogs.GET(sortable, app.Controllers.General.BlogController.GetSortableFields)
	}

	tickets := routerGroup.Group("/ticket")
	{
		tickets.GET(status, app.Controllers.General.TicketController.GetTicketStatuses)
		tickets.GET("/subject", app.Controllers.General.TicketController.GetTicketSubjects)
		tickets.GET(sortable, app.Controllers.General.TicketController.GetSortableFields)
	}

	bids := routerGroup.Group("/bid")
	{
		bids.GET(sortable, app.Controllers.General.BidController.GetSortableFields)
		bids.GET(status, app.Controllers.General.BidController.UserBidStatuses)
	}

	reports := routerGroup.Group("/report")
	{
		reports.GET(sortable, app.Controllers.General.ReportController.GetSortableFields)
		reports.GET(status, app.Controllers.General.ReportController.GetReportStatuses)
	}

	users := routerGroup.Group("/user")
	{
		users.GET(sortable, app.Controllers.General.UserController.GetSortableFields)
	}
}
