package broker

import (
	"fmt"

	subscription_messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	"google.golang.org/protobuf/proto"
)

type (
	Publisher interface {
		Publish(subject string, data []byte) error
	}

	natsBroker struct {
		conn Publisher
	}
)

func NewNatsBroker(connection Publisher) *natsBroker {
	return &natsBroker{
		conn: connection,
	}
}

func (b *natsBroker) CreateSubscription(sub deps.SubscriptionMsg) error {
	data, err := proto.Marshal(&subscription_messages.CreateSubscriptionMessage{})
	if err != nil {
		return fmt.Errorf("%w: failed to marshal proto message", err)
	}

	return b.conn.Publish("subscription.created", data)
}
