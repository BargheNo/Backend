package mqtt

type MQTTSubscription interface {
	SetupMQTTSubscriptions()
	RefreshMQTTSubscriptions()
}
