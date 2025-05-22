package serviceimpl

import (
	"github.com/BargheNo/Backend/internal/domain/mqtt"
)

type MonitoringService struct {
	mqttClient mqtt.Client
}

func NewMonitoringService(mqttClient mqtt.Client) *MonitoringService {
	return &MonitoringService{mqttClient: mqttClient}
}

func (s *MonitoringService) Test() {
	s.mqttClient.Subscribe("test")
}
