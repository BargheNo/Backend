package usecase

import monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"

type MonitoringService interface {
	HandleMessage(topic string, payload []byte)
	GetCustomerPanelStatus(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error)
	GetCustomerPanelHistory(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error)
	GetCustomerPanelEvent(listInfo monitoringdto.CustomerPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error)
	GetCorporationPanelStatus(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelStatusResponse, int64, error)
	GetCorporationPanelHistory(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelHistoryResponse, int64, error)
	GetCorporationPanelEvent(listInfo monitoringdto.CorporationPanelStatusListRequest) ([]monitoringdto.PanelEventResponse, int64, error)
}
