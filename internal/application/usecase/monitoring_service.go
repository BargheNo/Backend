package usecase

type MonitoringService interface {
	HandleMessage(topic string, payload []byte)
}
