package notification

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralNotificationController struct {
	constants           *bootstrap.Constants
	notificationService service.NotificationService
}

func NewGeneralNotificationController(
	constants *bootstrap.Constants,
	notificationService service.NotificationService,
) *GeneralNotificationController {
	return &GeneralNotificationController{
		constants:           constants,
		notificationService: notificationService,
	}
}

func (notificationController *GeneralNotificationController) GetContactTypes(ctx *gin.Context) {
	notificationTypes := notificationController.notificationService.GetNotificationsType()
	controller.Response(ctx, 200, "", notificationTypes)
}
