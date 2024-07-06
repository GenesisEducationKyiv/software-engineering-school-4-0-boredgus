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
		Subject() string
		Data() []byte
		Ack() error
		Nak() error
		Term() error
	}
	natsBroker struct {
		js              jetstream.JetStream
		eventConsumer   jetstream.Consumer
		commandConsumer jetstream.Consumer
		logger          config.Logger
	}
)

const CreationTimeout = 5 * time.Second

func createEventConsumer(js jetstream.JetStream) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CreationTimeout)
	defer cancel()

	eventStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "EVENTS",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"events.>"},
	})
	if err != nil {
		return nil, err
	}

	return eventStream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:          "notification-event-consumer",
		Durable:       "notification-event-consumer",
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
}

func createCommandConsumer(js jetstream.JetStream) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CreationTimeout)
	defer cancel()

	eventStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "COMMANDS",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"commands.>"},
	})
	if err != nil {
		return nil, err
	}

	return eventStream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:          "notification-command-consumer",
		Durable:       "notification-command-consumer",
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
}

func NewNatsBroker(js jetstream.JetStream, logger config.Logger) (*natsBroker, error) {

	eventConsumer, err := createEventConsumer(js)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS stream instance for events: %w", err)
	}

	commandConsumer, err := createCommandConsumer(js)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS stream consumer for commands: %w", err)
	}

	return &natsBroker{
		js:              js,
		eventConsumer:   eventConsumer,
		commandConsumer: commandConsumer,
		logger:          logger,
	}, nil
}

func (b *natsBroker) consume(consumer jetstream.Consumer, handler func(msg ConsumedMessage)) error {
	_, err := consumer.Consume(func(msg jetstream.Msg) {
		handler(msg)
	})

	if err != nil {
		return fmt.Errorf("failed to consume message from stream '%s': %w", consumer.CachedInfo().Stream, err)
	}

	return nil
}

func (b *natsBroker) ConsumeEvent(handler func(msg ConsumedMessage)) error {
	return b.consume(b.eventConsumer, handler)
}

func (b *natsBroker) ConsumeCommand(handler func(msg ConsumedMessage)) error {
	return b.consume(b.commandConsumer, handler)
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

func (b *natsBroker) ObjectStore(bucket string) (jetstream.ObjectStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CreationTimeout)
	defer cancel()

	store, err := b.js.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket:  bucket,
		Storage: jetstream.FileStorage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create object store for dispatches: %w", err)
	}

	return store, nil
}
