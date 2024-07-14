package broker

import (
	"errors"
	"fmt"

	"context"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	ConsumedMessage interface {
		// Subject returns a subject on which a message was published/received.
		Subject() string

		// Data returns the message body.
		Data() []byte
	}

	natsBroker struct {
		jetstream jetstream.JetStream
		consumer  jetstream.Consumer
		logger    config.Logger
	}
)

var SkippedMessageErr = errors.New("message is skipped")

const (
	CreationTimeout time.Duration = 5 * time.Second
	RedeliveryDelay time.Duration = 1 * time.Minute

	NotificationStreamName   string = "NOTIFICATIONS"
	NotificationConsumerName string = "notification-sender"
)

func CreateNotificationConsumer(js jetstream.JetStream) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CreationTimeout)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      NotificationStreamName,
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"commands.>", "events.subscription.>"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update %s stream: %w", NotificationStreamName, err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          NotificationConsumerName,
		Durable:       NotificationConsumerName,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update %s consumer: %w", NotificationConsumerName, err)
	}

	return consumer, nil
}

func NewNatsBroker(
	js jetstream.JetStream,
	consumer jetstream.Consumer,
	logger config.Logger,
) (*natsBroker, error) {
	return &natsBroker{
		jetstream: js,
		consumer:  consumer,
		logger:    logger,
	}, nil
}

func (b *natsBroker) ConsumeMessage(handler func(msg ConsumedMessage) error) error {
	_, err := b.consumer.Consume(func(msg jetstream.Msg) {
		err := handler(msg)

		if errors.Is(err, SkippedMessageErr) {
			b.logger.Infof("skipping message with subject %v ...", msg.Subject())

			return
		}

		if err != nil {
			b.logger.Errorf("failed to handle message: %v", err)

			err = msg.NakWithDelay(RedeliveryDelay)
			if err != nil {
				b.logger.Errorf("failed to negatively acknowledge message: %v", err)
			}

			return
		}

		err = msg.Ack()
		if err != nil {
			b.logger.Errorf("failed to acknowledge message: %v", err)
		}

		b.logger.Info("successfully handled message")
	})
	if err != nil {
		return fmt.Errorf("failed to consume message from stream '%s': %w", b.consumer.CachedInfo().Stream, err)
	}

	return nil
}

func (b *natsBroker) PublishAsync(subject string, payload []byte) error {
	pubAckFuture, err := b.jetstream.PublishAsync(
		subject,
		payload,
		jetstream.WithMsgID(uuid.NewString()),
	)
	go func() {
		msgID := pubAckFuture.Msg().Header.Get(jetstream.MsgIDHeader)
		select {
		case pubAck := <-pubAckFuture.Ok():

			b.logger.Infof("message with id '%s' asynchronously published to stream '%s'", msgID, pubAck.Stream)

		case err := <-pubAckFuture.Err():
			if err != nil {
				b.logger.Errorf("failed to publish message with '%s' to stream '%s': %v", err)
			}
		}
	}()

	return err
}

func (b *natsBroker) ObjectStore(bucket string) (jetstream.ObjectStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CreationTimeout)
	defer cancel()

	store, err := b.jetstream.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket:  bucket,
		Storage: jetstream.FileStorage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create object store: %w", err)
	}

	return store, nil
}
