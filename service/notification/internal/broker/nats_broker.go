package broker

import (
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

		// Ack tells the server that the message was successfully processed
		// and it can move on to the next message.
		Ack() error

		// NakWithDelay tells the server to redeliver the message after the given delay.
		NakWithDelay(delay time.Duration) error
	}

	natsBroker struct {
		jetstream jetstream.JetStream
		consumer  jetstream.Consumer
		logger    config.Logger
	}
)

const (
	CreationTimeout time.Duration = 5 * time.Second

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

func (b *natsBroker) ConsumeMessage(handler func(msg ConsumedMessage)) error {
	_, err := b.consumer.Consume(func(msg jetstream.Msg) {
		handler(msg)
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
