package mqtt

type Client interface {
	Subscribe(topic string, handler func(payload []byte))
	Disconnect()
}
