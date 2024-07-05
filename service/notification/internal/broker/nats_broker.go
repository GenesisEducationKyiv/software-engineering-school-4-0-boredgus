package broker

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type (
	natsBroker struct {
		conn *nats.Conn
	}
)

func NewNatsBroker(connection *nats.Conn) *natsBroker {
	return &natsBroker{
		conn: connection,
	}
}

func (b *natsBroker) Publish(subject string, data []byte) error {
	return b.conn.Publish(subject, data)
}

func (b *natsBroker) Consume(subject, queue string, handler func(msg *ConsumedMessage)) error {
	sub, err := b.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(NewConsumedMessage(msg))
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to queue '%s' with subject '%s': %w", queue, subject, err)
	}

	return sub.Drain()
}
