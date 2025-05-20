package messaging

import (
	"context"
	"sync"
)

var once sync.Once

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
	Consume(ctx context.Context, handler func(key, value []byte) error)
}
