package broker

import (
	"github.com/nats-io/nats.go/jetstream"
)

type (
	Broker interface {
		ConsumeMessage(handler func(msg ConsumedMessage) error) error
		PublishAsync(subject string, payload []byte) error
		ObjectStore(bucket string) (jetstream.ObjectStore, error)
	}
	Counter interface {
		Inc()
	}
	Metrics struct {
		TotalMessages              Counter
		SuccessfulyHandledMessages Counter
	}
	natsBrokerWithMetrics struct {
		broker                     Broker
		totalMessages              Counter
		successfulyHandledMessages Counter
	}
)

func NewNatsBrokerWithMetrics(broker Broker, metrics Metrics) *natsBrokerWithMetrics {
	return &natsBrokerWithMetrics{
		broker:                     broker,
		totalMessages:              metrics.TotalMessages,
		successfulyHandledMessages: metrics.SuccessfulyHandledMessages,
	}
}

func (b *natsBrokerWithMetrics) ConsumeMessage(handler func(msg ConsumedMessage) error) error {
	return b.broker.ConsumeMessage(func(msg ConsumedMessage) error {
		b.totalMessages.Inc()
		res := handler(msg)
		if res == nil {
			b.successfulyHandledMessages.Inc()
		}

		return res
	})
}

func (b *natsBrokerWithMetrics) PublishAsync(subject string, payload []byte) error {
	return b.broker.PublishAsync(subject, payload)
}

func (b *natsBrokerWithMetrics) ObjectStore(bucket string) (jetstream.ObjectStore, error) {
	return b.broker.ObjectStore(bucket)
}
