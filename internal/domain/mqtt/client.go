package mqtt

type Client interface {
	Subscribe(topic string, handler func(topic string, payload []byte))
	Disconnect()
}
