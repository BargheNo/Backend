package usecase

import monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"

type MonitoringService interface {
	HandleMessage(topic string, payload []byte)
	GetPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error)
	GetPanelHistory(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error)
	GetPanelEvent(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error)
}
