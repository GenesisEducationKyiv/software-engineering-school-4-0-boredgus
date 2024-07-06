package app

import (
	"context"
	"time"

	broker_msgs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Publisher interface {
		PublishAsync(subject string, payload []byte) error
	}
	Converter interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}
)

func DispatchInvoker(
	broker Publisher,
	converter Converter,
	logger config.Logger,
) func(d *entities.Dispatch) {

	return func(d *entities.Dispatch) {
		msg := broker_msgs.SendDispatchCommand{
			EventType: broker_msgs.EventType_SEND_DISPATCH,
			Timestamp: timestamppb.New(time.Now().UTC()),
			Data: &broker_msgs.Data{
				Emails:  d.Emails,
				BaseCcy: d.BaseCcy,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
		defer cancel()

		rates, err := converter.Convert(ctx, d.BaseCcy, d.TargetCcies)
		if err != nil {
			logger.Errorf("failed to get rates: %v", err)

			return
		}
		msg.Data.Rates = rates

		marshalled, err := proto.Marshal(&msg)
		if err != nil {
			logger.Errorf("failed to marshal SendDispatchCommand: %v", err)

			return
		}

		if err := broker.PublishAsync(SendDispatchCommand, marshalled); err != nil {
			logger.Errorf("failed to emit SendDispatch commands: %v", err)
		}
	}
}
