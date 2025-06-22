package notification

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralNotificationController struct {
	constants           *bootstrap.Constants
	notificationService port.NotificationService
}

func NewGeneralNotificationController(
	constants *bootstrap.Constants,
	notificationService port.NotificationService,
) *GeneralNotificationController {
	return &GeneralNotificationController{
		constants:           constants,
		notificationService: notificationService,
	}
}

func (notificationController *GeneralNotificationController) GetContactTypes(ctx *gin.Context) {
	notificationTypes, err := notificationController.notificationService.GetNotificationsType()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", notificationTypes)
}
