package broker

import (
	"context"

	"github.com/nats-io/nats.go"
)

type ConsumedMessage struct {
	msg *nats.Msg
}

func NewConsumedMessage(msg *nats.Msg) *ConsumedMessage {
	return &ConsumedMessage{msg: msg}
}

func (m *ConsumedMessage) Data() []byte {
	return m.msg.Data
}

func (m *ConsumedMessage) Ack(ctx context.Context) error {
	return m.msg.Ack(nats.Context(ctx))
}

func (m *ConsumedMessage) Nak(ctx context.Context) error {
	return m.msg.Nak(nats.Context(ctx))
}

func (m *ConsumedMessage) Terminate(ctx context.Context) error {
	return m.msg.Term(nats.Context(ctx))
}
