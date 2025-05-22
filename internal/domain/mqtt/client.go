package mqtt

type Client interface {
	Subscribe(topic string)
	Disconnect()
}
