package mocks

import "github.com/stretchr/testify/mock"

type BrokerMock struct {
	mock.Mock
}

func NewBrokerMock() *BrokerMock {
	return &BrokerMock{}
}

func (b *BrokerMock) Close() error {
	args := b.Called()
	return args.Error(0)
}

func (b *BrokerMock) ConsumeMessages(queue string, handler func([]byte) error) error {
	args := b.Called(queue, handler)
	return args.Error(0)
}

func (b *BrokerMock) PublishMessage(queue string, message interface{}) error {
	args := b.Called(queue, message)
	return args.Error(0)
}
