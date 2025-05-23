package service

type MonitoringService interface {
	HandleMessage(topic string, payload []byte)
}
