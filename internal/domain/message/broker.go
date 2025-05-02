package message

type Broker interface {
	Close() error
	ConsumeMessages(queue string, handler func([]byte) error) error
	PublishMessage(queue string, message interface{}) error
}
