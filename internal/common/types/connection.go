package types

import (
	"context"
)

type Connection interface {
	Close()
	IsConnected() bool
}
type ProducerConnection interface {
	Connection
	Produce(message string) error
}

type ConsumerConnection interface {
	Connection
	StartConsumer(ctx context.Context, ms MessageProcessor)
}

type MessageProcessor interface {
	ProcessMessage(message string) error
}
