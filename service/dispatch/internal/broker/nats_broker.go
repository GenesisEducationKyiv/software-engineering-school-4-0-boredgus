package broker

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
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

func NewNatsBroker(conn *nats.Conn, logger config.Logger, onError func(error, string)) *natsBroker {
	js, err := jetstream.New(
		conn,
		jetstream.WithPublishAsyncErrHandler(PublishAsyncErrHandler(logger)),
	)
	onError(err, "failed to create NATS Jetstream instance")

	return &natsBroker{
		js:     js,
		logger: logger,
	}
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