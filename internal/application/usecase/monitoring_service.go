package usecase

import monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"

type MonitoringService interface {
	HandleMessage(topic string, payload []byte)
	GetPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.CustomerPanelStatusResponse, int64, error)
}
