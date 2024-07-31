package broker

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type natsBroker struct {
	js     jetstream.JetStream
	logger config.Logger
}

func PublishAsyncErrHandler(logger config.Logger) jetstream.MsgErrHandler {
	return func(js jetstream.JetStream, m *nats.Msg, err error) {
		logger.Errorf("handler: failed to publish message '%s' asynchronously: %v", string(m.Data), err)
	}
}

func NewNatsBroker(conn *nats.Conn, logger config.Logger) (*natsBroker, error) {
	js, err := jetstream.New(
		conn,
		jetstream.WithPublishAsyncErrHandler(PublishAsyncErrHandler(logger)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create jestream: %w", err)
	}

	return &natsBroker{
		js:     js,
		logger: logger,
	}, nil
}

func (b *natsBroker) PublishAsync(subject string, payload []byte) error {
	pubAckFuture, err := b.js.PublishAsync(
		subject,
		payload,
		jetstream.WithMsgID(uuid.NewString()),
	)
	go func() {
		select {
		case pubAck := <-pubAckFuture.Ok():
			b.logger.Infof("message %s asynchronously published to stream '%s'", pubAckFuture.Msg().Data, pubAck.Stream)

		case err := <-pubAckFuture.Err():
			if err != nil {
				b.logger.Errorf("failed to publish message to stream: %v", err)
			}
		}
	}()

	return err
}
