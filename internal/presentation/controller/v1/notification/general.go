package notification

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralNotificationController struct {
	constants           *bootstrap.Constants
	notificationService usecase.NotificationService
}

func NewGeneralNotificationController(
	constants *bootstrap.Constants,
	notificationService usecase.NotificationService,
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

func (notificationController *GeneralNotificationController) GetSortableFields(ctx *gin.Context) {
	columns := notificationController.notificationService.GetNotificationSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
